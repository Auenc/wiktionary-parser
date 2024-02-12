package wikitionaryparser

type Detail interface {
	Name() string
}

func detailsFromSections(language string, sections map[string]string) ([]Detail, error) {
	details := make([]Detail, 0)

	for sectionName, section := range sections {
		detail, err := detailFromSection(language, sectionName, section)
		if err != nil {
			return details, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func detailFromSection(language, sectionName, section string) (Detail, error) {
	var detail Detail
	// allow us to overwrite default details with language specific details
	switch language {

	}
	switch sectionName {
	// case "romanization":
	// 	return Romanization{}
	}
	return detail, nil
}
