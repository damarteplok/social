package main

type FormDataBikinSomething struct {
	Text1 string `json:"text1" validate:"required"`
	EmailTest string `json:"email_test" validate:"required"`
	TextAreaTest string `json:"text_area_test,omitempty"`
	NumberTest float64 `json:"number_test,omitempty"`
	SelectTest string `json:"select_test,omitempty"`
}

	
