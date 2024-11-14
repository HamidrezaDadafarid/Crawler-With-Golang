package rate_limiter

func getLimitByRole(role string) int {
    switch role {
    case "admin":
        return 100
    case "super_user":
        return 50
    default:
        return 20
    }
}
