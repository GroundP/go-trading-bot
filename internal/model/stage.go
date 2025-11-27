// Package model
package model

type StageNumber int

const (
	STAGE_0 StageNumber = iota // 초기값
	STAGE_1                    // 안정 상승기, 단/중/장 배치, 매수 진입
	STAGE_2                    // 데드크로스, 중/단/장 배치
	STAGE_3                    // 본격 하락기, 중/장/단 배치, 매수 청산
	STAGE_4                    // 안정 하락기, 장/중/단 배치, 매도 진입
	STAGE_5                    // 골든크로스, 장/단/중 배치
	STAGE_6                    // 본격 상승기, 단/장/중 배치, 매도 청산
)

func (s StageNumber) String() string {
	return [...]string{"STAGE_0", "STAGE_1", "STAGE_2", "STAGE_3", "STAGE_4", "STAGE_5", "STAGE_6"}[s]
}

type StageDir string

const (
	STAGE_DIR_NONE     StageDir = "NONE"     // 초기값
	STAGE_DIR_MAINTAIN StageDir = "MAINTAIN" // 유지
	STAGE_DIR_NORMAL   StageDir = "NORMAL"   // 정상 방향
	STAGE_DIR_REVERSE  StageDir = "REVERSE"  // 역방향
)

type Stage struct {
	StageNumber StageNumber
	StageDir    StageDir
	Description string
}
