package users

// Helper function to handle nil pointers for string values.
func defaultString(s *string) string {
    if s != nil {
        return *s
    }
    return ""
}

// Helper function to handle nil pointers for int values.
func defaultInt(i *int) int {
    if i != nil {
        return *i
    }
    return 0
}
