package services

import (
	"bufio"
	"fmt"
	"order_food_online/internal/cache"
	"order_food_online/internal/models"
	"os"
	"strings"
	"time"
)

type PromoCodeService struct {
	cache cache.PromoCodeCache
}

func NewPromoCodeService(cache cache.PromoCodeCache) *PromoCodeService {
	return &PromoCodeService{cache: cache}
}

func (s *PromoCodeService) ValidatePromo(code string) (bool, error) {
	// Validate code length
	if len(code) < 8 || len(code) > 10 {
		return false, nil
	}

	// Check cache for the promo code
	cachedPromo, err := s.cache.GetPromoCode(code)
	if err == nil && cachedPromo != nil {
		return cachedPromo.IsValid, nil
	}

	// Validate against files
	isValid, err := validatePromoCodeFromFiles(code)
	if err != nil {
		return false, err
	}

	// Cache the result
	ttl := 24 * time.Hour // Cache the result for 1 day
	_ = s.cache.SetPromoCode(code, &models.PromoCode{Code: code, IsValid: isValid}, ttl)

	return isValid, nil
}

func validatePromoCodeFromFiles(code string) (bool, error) {
	files := []string{"couponbase1.txt", "couponbase2.txt", "couponbase3.txt"}
	matchCount := 0

	// Check each file for the promo code
	for _, file := range files {
		found, err := searchInFile(fmt.Sprintf("%s/%s", os.Getenv("COUPON_DIR"), file), code)
		if err != nil {
			return false, err
		}
		if found {
			matchCount++
		}
		if matchCount >= 2 {
			return true, nil
		}
	}

	return false, nil
}

func searchInFile(filePath, code string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == code {
			return true, nil
		}
	}
	return false, scanner.Err()
}
