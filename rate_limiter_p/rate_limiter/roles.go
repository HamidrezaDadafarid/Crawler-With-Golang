
package rate_limiter

func getLimitByRole(role string) int {
    switch role {
    case "admin":
        return 100
    case "premium_user":
        return 50
    default:
        return 20
    }
}
