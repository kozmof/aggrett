package aggrett

// FilterByTag keeps factors whose tags match the provided list.
func FilterByTag(sequence []SeqFactor, tags []string) []SeqFactor {
	tagSet := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}

	result := make([]SeqFactor, 0, len(sequence))
	for _, f := range sequence {
		if _, ok := tagSet[f.Tag]; ok {
			result = append(result, f)
		}
	}
	return result
}

// ExtractTags returns unique tags in first-seen order.
func ExtractTags(sequence []SeqFactor) []string {
	seen := make(map[string]struct{})
	tags := make([]string, 0)
	for _, f := range sequence {
		if _, ok := seen[f.Tag]; ok {
			continue
		}
		seen[f.Tag] = struct{}{}
		tags = append(tags, f.Tag)
	}
	return tags
}

// GroupByTag groups factors by tag.
func GroupByTag(sequence []SeqFactor) map[string][]SeqFactor {
	groups := make(map[string][]SeqFactor)
	for _, f := range sequence {
		groups[f.Tag] = append(groups[f.Tag], f)
	}
	return groups
}

// ExcludeByTag removes factors whose tags match the provided list.
func ExcludeByTag(sequence []SeqFactor, tags []string) []SeqFactor {
	tagSet := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}

	result := make([]SeqFactor, 0, len(sequence))
	for _, f := range sequence {
		if _, ok := tagSet[f.Tag]; !ok {
			result = append(result, f)
		}
	}
	return result
}

// RemoveByTag is an alias of ExcludeByTag.
func RemoveByTag(sequence []SeqFactor, tags []string) []SeqFactor {
	return ExcludeByTag(sequence, tags)
}

// RenameTag renames factors from oldTag to newTag.
func RenameTag(sequence []SeqFactor, oldTag, newTag string) []SeqFactor {
	result := make([]SeqFactor, 0, len(sequence))
	for _, f := range sequence {
		if f.Tag == oldTag {
			f.Tag = newTag
		}
		result = append(result, f)
	}
	return result
}

// AccumulateByTag accumulates only factors with the given tag.
func AccumulateByTag(sequence []SeqFactor, baseValue float64, tag string) []AccumCore {
	return AccumulateSequence(FilterByTag(sequence, []string{tag}), baseValue)
}
