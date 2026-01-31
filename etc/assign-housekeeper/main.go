package assignhk

import (
	"context"
	"errors"
	"sort"
	"time"
)

// ---------- Domain models (ย่อ/ตัวอย่าง) ----------
type Order struct {
	ID        int64
	BranchID  int64
	SlotID    int64     // time-slot (เช่น 09:00-10:00)
	ShiftID   int64     // กะเช้า/บ่าย/ดึก หรือ bucket อื่นๆ ที่ธุรกิจกำหนด
	StartTime time.Time // เวลาเริ่มของ slot (ถ้าใช้ time จริง)
	EndTime   time.Time // เวลาจบของ slot
}

type Housekeeper struct {
	ID        string
	Name      string
	BranchID  int64
	QueueTurn int64 // จำนวน round-robin/คิว (ตัวเลขเพิ่มตามเวลาถูกเลือก)
}

// Assignment ปัจจุบันที่ “ล็อก” อยู่ในระบบแล้ว
type Assignment struct {
	OrderID       int64
	HousekeeperID int64
	SlotID        int64
	ShiftID       int64
}

// ---------- Repository ports (ให้ต่อ GORM ได้) ----------
type Repository interface {
	// assignments ที่ชนกับช่วงเวลานี้ (หรือ slot นี้) ของแม่บ้านแต่ละคน
	GetAssignmentsInWindow(ctx context.Context, hkIDs []string, start, end time.Time) (map[string][]Assignment, error)

	// assignments ใน slot เดียวกันทั้งหมดของแม่บ้าน (เอาไว้เช็ค “มีงานใน slot นี้แล้วหรือยัง”)
	GetAssignmentsInSlot(ctx context.Context, hkIDs []string, slotID int64) (map[string][]Assignment, error)
}

// ---------- Service ----------
type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) pickAssignee(ctx context.Context, order Order, candidates []Housekeeper, useRoundRobin bool) (*Housekeeper, error) {
	if len(candidates) == 0 {
		return nil, errors.New("no candidates")
	}

	hkIDs := make([]string, 0, len(candidates))
	for _, hk := range candidates {
		hkIDs = append(hkIDs, hk.ID)
	}

	// 0) ถ้ามีคนเดียว ให้ assign ได้เลย
	if len(candidates) == 1 {
		return &candidates[0], nil
	}

	// 1) เช็ค “กะซ้ำ/ทับช่วงเวลา (Shift/TimeSlot overlap?)”)
	//    - ถ้า “ไม่ซ้ำ” → assign ได้เลย (จบตาม flow)
	//    - ถ้า “ซ้ำ”    → ไปขั้นตอนที่ 2
	winMap, err := s.repo.GetAssignmentsInWindow(ctx, hkIDs, order.StartTime, order.EndTime)
	if err != nil {
		return nil, err
	}
	notOverlapped := filterByNoOverlap(candidates, winMap)
	if len(notOverlapped) > 0 {
		chosen := notOverlapped[0] // เลือกคนแรกที่ “ไม่ซ้ำกะ”
		return &chosen, nil
	}

	// 2) เช็ค “งานซ้ำใน slot เดียวกันไหม”
	//    - ถ้ายัง “ไม่มีงานใน slot” → assign ให้คนนั้น
	//    - ถ้ามีงานใน slot อยู่แล้ว → ไปขั้นตอนที่ 3
	slotMap, err := s.repo.GetAssignmentsInSlot(ctx, hkIDs, order.SlotID)
	if err != nil {
		return nil, err
	}
	noWorkInThisSlot := filterByNoWorkInSlot(candidates, slotMap)
	if len(noWorkInThisSlot) > 0 {
		// กระจายงาน: เลือกคนที่ยังว่างใน slot นี้ก่อน (กัน error case งานกองอยู่ที่คนเดียว)
		chosen := noWorkInThisSlot[0]
		return &chosen, nil
	}

	// 3) ทุกคนใน slot นี้ “มีงานแล้ว”
	//    - ถ้ายังมีคนอื่น “ว่าง” (ในความหมายธุรกิจอื่น ๆ ที่อยากเช็คเพิ่ม เช่น workload น้อยกว่า)
	//      ก็เลือกคนนั้น (ตัวอย่างนี้ถือว่าไม่มีคนว่างแล้วเพราะข้อ 2 คัดออกไปแล้ว)
	//    - มิฉะนั้น ใช้ลำดับที่กำหนดไว้ (Queue)
	if useRoundRobin {
		// เลือกตามคิว (QueueTurn น้อยสุดก่อน)
		rr := append([]Housekeeper(nil), candidates...)
		sort.SliceStable(rr, func(i, j int) bool {
			return rr[i].QueueTurn < rr[j].QueueTurn
		})
		chosen := rr[0]
		return &chosen, nil
	}

	// ถ้าไม่ใช้ RR ให้ fallback ตาม Priority
	pp := append([]Housekeeper(nil), candidates...)
	chosen := pp[0]
	return &chosen, nil
}

// ---------- Helpers ----------

// แม่บ้านที่ “ไม่ทับช่วงเวลา/กะ” กับงานที่มีอยู่
func filterByNoOverlap(candidates []Housekeeper, winMap map[string][]Assignment) []Housekeeper {
	out := make([]Housekeeper, 0, len(candidates))
	for _, hk := range candidates {
		// ในตัวอย่างนี้ winMap ถูก query ด้วย window = slot time แล้ว
		// ถ้าไม่มี assignment ในช่วงเวลานี้ => ไม่ overlap
		if len(winMap[hk.ID]) == 0 {
			out = append(out, hk)
		}
	}
	return out
}

// แม่บ้านที่ “ยังไม่มีงานใน slot นี้”
func filterByNoWorkInSlot(candidates []Housekeeper, slotMap map[string][]Assignment) []Housekeeper {
	out := make([]Housekeeper, 0, len(candidates))
	for _, hk := range candidates {
		if len(slotMap[hk.ID]) == 0 {
			out = append(out, hk)
		}
	}
	return out
}
