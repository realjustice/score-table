package main

import (
	"fmt"
	. "score_table"
)

func main() {
	scoreTable := NewScoreTable(WithSOS(), WithSOSOS())
	// players
	// 1.KunneStéphan
	// 2.BuffardEmmanuel
	// 3.GauthierHenri
	// 4.Papazoglou Benjamin
	// 5.CancePhilippe
	// 6.KaradabanDenis
	// 7.ParcoitDavid
	// 8.Nguyen Huu_Phuoc
	// 9.TecchioPierre
	// 10.VannierRémi

	// Num  Name             NBW  SOS  SOSOS
	//1     KunneStéphan      3   14   64
	//2     BuffardEmmanuel   3   13   65
	//3     GauthierHenri     3   13   65
	//4     Papazoglou        3   13   64
	//5     CancePhilippe     3   13   63
	//6     KaradabanDenis    2   13   61
	//7     ParcoitDavid      2   12   61
	//8     Nguyen Huu_Phuoc  2   12   60
	//9     TecchioPierre     2   11   61
	//10    VannierRémi       2   11   61

	// KunneStéphan
	scoreTable.RecordResult(1, 1, 5, WHITE_WIN)
	scoreTable.RecordResult(2, 1, 2, BLACK_WIN)
	scoreTable.RecordResult(3, 1, 7, BLACK_WIN)
	scoreTable.RecordResult(4, 1, 4, BLACK_WIN)
	scoreTable.RecordResult(5, 1, 3, WHITE_WIN)
	// BuffardEmmanuel
	scoreTable.RecordResult(1, 2, 6, BLACK_WIN)
	scoreTable.RecordResult(3, 2, 8, BLACK_WIN)
	scoreTable.RecordResult(4, 2, 5, BLACK_WIN)
	scoreTable.RecordResult(5, 2, 4, WHITE_WIN)
	// GauthierHenri
	scoreTable.RecordResult(1, 3, 4, BLACK_WIN)
	scoreTable.RecordResult(2, 3, 7, BLACK_WIN)
	scoreTable.RecordResult(3, 3, 5, WHITE_WIN)
	scoreTable.RecordResult(4, 3, 6, WHITE_WIN)

	// Papazoglou Benjamin
	scoreTable.RecordResult(2, 4, 9, BLACK_WIN)
	scoreTable.RecordResult(3, 4, 6, BLACK_WIN)

	// CancePhilippe
	scoreTable.RecordResult(2, 5, 8, BLACK_WIN)
	scoreTable.RecordResult(5, 5, 10, WHITE_WIN)

	// KaradabanDenis
	scoreTable.RecordResult(2, 6, 10, BLACK_WIN)
	scoreTable.RecordResult(5, 6, 9, WHITE_WIN)

	// ParcoitDavid
	scoreTable.RecordResult(1, 7, 9, BLACK_WIN)
	scoreTable.RecordResult(4, 7, 10, WHITE_WIN)
	scoreTable.RecordResult(5, 7, 8, BLACK_WIN)

	// Nguyen Huu_Phuoc
	scoreTable.RecordResult(1, 8, 10, BLACK_WIN)
	scoreTable.RecordResult(4, 8, 9, BLACK_WIN)

	// TecchioPierre
	scoreTable.RecordResult(3, 9, 10, BLACK_WIN)

	scores := scoreTable.GetScoreTableByRound(5)

	scoreTable.GetPlayerScoreByRound(10, 5)
	for _, score := range scores {
		fmt.Printf("%+v\n", score)
	}
}
