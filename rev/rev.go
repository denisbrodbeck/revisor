// Package rev contains the library for static asset revisioning.
package rev

import (
	"log"
	"path/filepath"
)

func Revision(webroot, baseURL string, atomTypes, moleculeTypes, pageTypes []string) error {
	root, err := filepath.Abs(webroot)
	if err != nil {
		return err
	}
	log.Println("Web Root:", root)
	log.Println("Base URL:", baseURL)
	files, err := Find(root, baseURL)
	if err != nil {
		return err
	}
	// revision all binary assets
	atoms := filterFilesBySuffix(files, atomTypes...)
	if err := revisionAtoms(atoms); err != nil {
		return err
	}
	log.Printf("Total atoms: %d", len(atoms))
	// revision all molecule files
	molecules := filterFilesBySuffix(files, moleculeTypes...)
	if err := revisionMolecules(molecules, atoms); err != nil {
		return err
	}
	log.Printf("Total molecules: %d", len(molecules))
	// update all page files
	pages := filterFilesBySuffix(files, pageTypes...)
	if err := revisionPages(pages, molecules, atoms); err != nil {
		return err
	}
	log.Printf("Total pages: %d", len(pages))
	log.Printf("Total revisioned assets: %d", len(atoms)+len(molecules))
	return nil
}

func revisionAtoms(atoms Files) error {
	for i := range atoms {
		if err := atoms[i].Revision(); err != nil {
			return err
		}
	}
	return nil
}

func revisionMolecules(molecules, atoms Files) error {
	for i := range molecules {
		if err := molecules[i].ReplaceReferences(atoms); err != nil {
			return err
		}
		if err := molecules[i].Revision(); err != nil {
			return err
		}
	}
	return nil
}

func revisionPages(pages, molecules, atoms Files) error {
	refs := Files{}
	refs = append(refs, atoms...)
	refs = append(refs, molecules...)

	for i := range pages {
		if err := pages[i].ReplaceReferences(refs); err != nil {
			return err
		}
	}
	return nil
}
