package rest

type (
	NewFlashcardRequest struct {
		NativeLanguage string `json:"native_lang"`
		TargetLang     string `json:"target_lang"`
		Content        struct {
			WordInNative  string   `json:"word_in_native"`
			WordInTarget  string   `json:"word_in_target"`
			UsageExamples []string `json:"usage"`
		} `json:"content"`
	}

	NewFlashcardResponse struct {
		Errors []string `json:"errors,omitempty"` // May be empty in OK
	}

	GetFlashcardRequest struct {
	}
)

// TODO grammar cards
