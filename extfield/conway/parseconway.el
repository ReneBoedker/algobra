;; The following Emacs commands turn CPimport.txt into the Go case statements
(progn
  (save-mark-and-excursion
   (goto-char 0)
   (replace-string "allConwayPolynomials := ["
				   "	switch {")
   (replace-regexp "\\[\\([0-9]+\\),\\([0-9]+\\),\\[\\([^]]*\\)\\]\\],"
				   "	case char==\\1 && extDeg==\\2:
		return []uint{\\3},nil")
   (replace-string "0];"
				   "	default:
		return nil, errors.New(
			op, errors.InputValue,
			\"No Conway polynomial was found for characteristic %d and extension degree %d\",
			char,extDeg,
		)
	}")))
