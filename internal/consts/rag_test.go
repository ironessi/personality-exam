package consts

import (
	"testing"
)

// 42 个测评里各策略的数量。这个数字是附录 B 的校验行，改语料或改策略必须同步。
const (
	wantFull       = 38
	wantCommonOnly = 2
	wantDisabled   = 2
	wantTotal      = 42
)

// TestRegistryCoversAllTestTypes 保证注册表不漏掉任何一个 consts.go 里的测评。
// 漏登记的后果不是报错，而是**静默不增强**——最难发现的那种问题。
func TestRegistryCoversAllTestTypes(t *testing.T) {
	for code := 1; code <= wantTotal; code++ {
		if _, ok := TestRegistry[code]; !ok {
			t.Errorf("测评编码 %d 未登记进 TestRegistry", code)
		}
	}
	if len(TestRegistry) != wantTotal {
		t.Errorf("TestRegistry 有 %d 条，期望 %d 条", len(TestRegistry), wantTotal)
	}
}

func TestRegistryKeysAreUnique(t *testing.T) {
	seen := make(map[string]int, len(TestRegistry))
	for code := 1; code <= wantTotal; code++ {
		m, ok := TestRegistry[code]
		if !ok {
			continue
		}
		if m.Key == "" {
			t.Errorf("编码 %d 的 Key 为空", code)
		}
		if prev, dup := seen[m.Key]; dup {
			t.Errorf("Key %q 重复：编码 %d 与 %d", m.Key, prev, code)
		}
		seen[m.Key] = code
	}
}

// TestEveryMetaIsSelfConsistent 防止「加了条目但填错了」这类复制粘贴错误。
func TestEveryMetaIsSelfConsistent(t *testing.T) {
	for code := 1; code <= wantTotal; code++ {
		m, ok := TestRegistry[code]
		if !ok {
			continue
		}
		if m.Code != code {
			t.Errorf("TestRegistry[%d].Code = %d，键值应一致", code, m.Code)
		}
		if m.Name == "" {
			t.Errorf("编码 %d 的 Name 为空", code)
		}
		switch m.Policy {
		case PolicyFull:
			if m.Domain == "" || m.Domain == DomainCommon {
				t.Errorf("编码 %d 是 full，Domain 必须是具体构念域，实为 %q", code, m.Domain)
			}
		case PolicyCommonOnly:
			if m.Domain == "" {
				t.Errorf("编码 %d 是 common_only，Domain 仍应登记（只是不参与检索）", code)
			}
		case PolicyDisabled:
			if m.Domain != "" {
				t.Errorf("编码 %d 是 disabled，不应有 Domain，实为 %q", code, m.Domain)
			}
		default:
			t.Errorf("编码 %d 的 Policy %q 不是合法值", code, m.Policy)
		}
		if m.Domain != "" && !IsValidDomain(m.Domain) {
			t.Errorf("编码 %d 的 Domain %q 不在 AllDomains 里", code, m.Domain)
		}
	}
}

// TestPolicyDistribution 钉死 38/2/2——附录 B 的校验行。
// 改这里意味着改了「多少测评会向用户输出个性化解读」，属于安全影响，必须是有意为之。
func TestPolicyDistribution(t *testing.T) {
	var full, commonOnly, disabled int
	for _, m := range TestRegistry {
		switch m.Policy {
		case PolicyFull:
			full++
		case PolicyCommonOnly:
			commonOnly++
		case PolicyDisabled:
			disabled++
		}
	}
	if full != wantFull || commonOnly != wantCommonOnly || disabled != wantDisabled {
		t.Errorf(
			"策略分布 = full:%d common_only:%d disabled:%d，期望 %d/%d/%d",
			full, commonOnly, disabled, wantFull, wantCommonOnly, wantDisabled,
		)
	}
}

func TestRetrievalDomains(t *testing.T) {
	// full：本域 + common，顺序固定（先本域后 common）
	got := RetrievalDomains(MBTI)
	want := []string{"typology", "common"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("RetrievalDomains(MBTI) = %v，期望 %v", got, want)
	}

	// common_only：只有 common，本域不参与检索
	if got := RetrievalDomains(StanfordBinet); len(got) != 1 || got[0] != "common" {
		t.Errorf("RetrievalDomains(StanfordBinet) = %v，期望 [common]", got)
	}
	if got := RetrievalDomains(Rorschach); len(got) != 1 || got[0] != "common" {
		t.Errorf("RetrievalDomains(Rorschach) = %v，期望 [common]", got)
	}

	// disabled / 未登记：一律 nil，调用方据此跳过增强
	for _, code := range []int{Zootopia, SBTI, 9999} {
		if got := RetrievalDomains(code); got != nil {
			t.Errorf("RetrievalDomains(%d) = %v，期望 nil", code, got)
		}
	}
}

// TestFullPolicyNeverDropsCommon 是防止「检索范围漏掉通用层」的兜底：
// common 里有免责声明、求助渠道、量表局限三块所有报告都需要的语料。
func TestFullPolicyNeverDropsCommon(t *testing.T) {
	for code, m := range TestRegistry {
		if m.Policy != PolicyFull {
			continue
		}
		domains := RetrievalDomains(code)
		found := false
		for _, d := range domains {
			if d == string(DomainCommon) {
				found = true
			}
		}
		if !found {
			t.Errorf("编码 %d（full）的检索范围 %v 不含 common", code, domains)
		}
	}
}

func TestCanEnhanceAndNeedsReferral(t *testing.T) {
	if !CanEnhance(MBTI) {
		t.Error("MBTI 应当可增强")
	}
	if CanEnhance(Zootopia) || CanEnhance(SBTI) {
		t.Error("娱乐测评不应接入增强")
	}
	if !CanEnhance(StanfordBinet) {
		t.Error("common_only 的测评仍要调用增强（只是只检索 common）")
	}
	if CanEnhance(9999) {
		t.Error("未登记测评不应被判为可增强")
	}

	if !NeedsReferralForced(StanfordBinet) || !NeedsReferralForced(Rorschach) {
		t.Error("common_only 的两个测评必须强制 needs_referral")
	}
	if NeedsReferralForced(MBTI) {
		t.Error("full 的测评不应被强制 needs_referral")
	}
}

// TestDomainTestCountsMatchesAppendixB 与 doc/RAG增强系统开发文档.md 附录 B 对账。
// 注意 typology 是 5（4 个 full + 罗夏），ability 是 2（MHLAQ 走 full、斯坦福走 common_only）。
func TestDomainTestCountsMatchesAppendixB(t *testing.T) {
	want := map[Domain]int{
		DomainSymptom:   11,
		DomainTrait:     11,
		DomainTypology:  5,
		DomainSocial:    6,
		DomainWellbeing: 4,
		DomainCareer:    1,
		DomainAbility:   2,
		DomainCommon:    0,
	}
	got := DomainTestCounts()
	for domain, n := range want {
		if got[domain] != n {
			t.Errorf("域 %s 登记了 %d 个测评，期望 %d", domain, got[domain], n)
		}
	}

	// 登记了构念域的条目 + disabled（无域）应恰好等于 42
	sum := 0
	for _, n := range got {
		sum += n
	}
	if sum != wantTotal-wantDisabled {
		t.Errorf("有域的条目共 %d 个，期望 %d（42 − 2 个 disabled）", sum, wantTotal-wantDisabled)
	}
}

func TestRegisteredCodesIsCompleteAndSorted(t *testing.T) {
	codes := RegisteredCodes()
	if len(codes) != wantTotal {
		t.Fatalf("RegisteredCodes 返回 %d 个，期望 %d 个", len(codes), wantTotal)
	}
	for i, code := range codes {
		if code != i+1 {
			t.Fatalf("第 %d 个编码是 %d，说明有缺口或乱序：%v", i, code, codes)
		}
	}
}

// TestAllDomainsHasNoDuplicate 防止 AllDomains 被复制粘贴改坏——
// 它与 data/raw 目录名、Python 侧 KB_NAMESPACES 三处必须完全一致。
func TestAllDomainsHasNoDuplicate(t *testing.T) {
	if len(AllDomains) != 8 {
		t.Fatalf("AllDomains 有 %d 个域，期望 8 个", len(AllDomains))
	}
	seen := map[Domain]bool{}
	for _, d := range AllDomains {
		if seen[d] {
			t.Errorf("域 %q 重复", d)
		}
		seen[d] = true
	}
	if !seen[DomainCommon] {
		t.Error("AllDomains 必须包含 common")
	}
}
