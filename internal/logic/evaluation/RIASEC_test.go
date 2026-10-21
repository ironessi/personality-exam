package evaluation

import "testing"

// TestRIASECTieBreakIsDeterministic 盯的是「同一份作答每次得到同一个类型」。
// 原实现 range map 找最大值，Go 的 map 迭代顺序随机 —— 54/54 并列时
// dominant_type 会在「社会型(S)」和「研究型(I)」之间跳，用户重测一次结果就变。
func TestRIASECTieBreakIsDeterministic(t *testing.T) {
	const runs = 50

	first, second := determineRIASECTypes(36, 54, 36, 54, 36, 54)
	if first != "研究型(I)" || second != "社会型(S)" {
		t.Fatalf("并列时的定序不对：dominant=%q secondary=%q", first, second)
	}

	for i := 0; i < runs; i++ {
		d, s := determineRIASECTypes(36, 54, 36, 54, 36, 54)
		if d != first || s != second {
			t.Fatalf("第 %d 次得到 (%q,%q)，与首次 (%q,%q) 不一致", i, d, s, first, second)
		}
	}
}

// TestRIASECNoTie 确认非并列时仍然按分数取。
func TestRIASECNoTie(t *testing.T) {
	d, s := determineRIASECTypes(10, 20, 30, 40, 50, 60)
	if d != "常规型(C)" || s != "企业型(E)" {
		t.Fatalf("dominant=%q secondary=%q", d, s)
	}
}
