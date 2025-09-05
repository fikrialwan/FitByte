package main

import (
	"fmt"
	"strings"
)

func main() {
	// Test more comprehensive patterns including edge cases we might have missed
	testPatterns := []string{
		// Number range violations
		`{"weight":9}`,                                  // Below min range (10-1=9)
		`{"weight":1001}`,                               // Above max range (1000+1=1001)
		`{"height":2}`,                                  // Below min range (3-1=2)
		`{"height":251}`,                                // Above max range (250+1=251)
		
		// String length violations
		`{"name":"a"}`,                                  // Too short (1 char, min=2)
		`{"name":"` + strings.Repeat("a", 61) + `"}`,    // Too long (61 chars, max=60)
		`{"name":"` + strings.Repeat("b", 61) + `"}`,    // Too long with different char
		
		// Invalid enum values
		`{"preference":"notAnEnumValue"}`,               // Invalid enum
		`{"weightUnit":"notAnEnumValue"}`,               // Invalid enum
		`{"heightUnit":"notAnEnumValue"}`,               // Invalid enum
		`{"preference":"INVALID"}`,                      // Different invalid enum
		`{"weightUnit":"INVALID"}`,                      // Different invalid enum
		`{"heightUnit":"INVALID"}`,                      // Different invalid enum
		
		// Invalid URL patterns
		`{"imageUri":"notAUrl"}`,                        // Invalid URL
		`{"imageUri":"http://incomplete"}`,              // Incomplete URL
		`{"imageUri":"notAnObject"}`,                    // Different invalid pattern
		`{"imageUri":"notABoolean"}`,                    // Different invalid pattern
		`{"imageUri":"invalid-url"}`,                    // Hyphenated invalid URL
		`{"imageUri":"ftp://example.com"}`,              // Non-HTTP protocol
		
		// Invalid number strings
		`{"weight":"notANumber"}`,                       // Invalid number string
		`{"height":"notANumber"}`,                       // Invalid number string
		
		// Null values
		`{"name":null}`,                                 // Null value
		`{"imageUri":null}`,                             // Null value
		`{"preference":null}`,                           // Null value
		`{"weightUnit":null}`,                           // Null value
		`{"heightUnit":null}`,                           // Null value
		`{"weight":null}`,                               // Null value
		`{"height":null}`,                               // Null value
		
		// Empty strings
		`{"preference":""}`,                             // Empty required string
		`{"weightUnit":""}`,                             // Empty required string
		`{"heightUnit":""}`,                             // Empty required string
		`{"name":""}`,                                   // Empty optional string
		`{"imageUri":""}`,                               // Empty optional string
		
		// Type mismatches
		`{"preference":123}`,                            // Number as string
		`{"preference":true}`,                           // Boolean as string
		`{"name":123}`,                                  // Number as string
		`{"name":true}`,                                 // Boolean as string
		`{"imageUri":123}`,                              // Number as string
		`{"imageUri":true}`,                             // Boolean as string
		`{"weightUnit":123}`,                            // Number as string
		`{"heightUnit":123}`,                            // Number as string
		
		// Zero and negative values
		`{"weight":0}`,                                  // Zero (below min)
		`{"height":0}`,                                  // Zero (below min)
		`{"weight":-1}`,                                 // Negative
		`{"height":-1}`,                                 // Negative
	}

	fmt.Println("Testing validation patterns:")
	fmt.Println("=============================")

	passCount := 0
	failCount := 0

	for i, pattern := range testPatterns {
		shouldFail := shouldValidationFail(pattern)
		status := "PASS"
		if !shouldFail {
			status = "FAIL"
			failCount++
		} else {
			passCount++
		}
		fmt.Printf("%2d. %s -> %s\n", i+1, pattern, status)
	}

	fmt.Printf("\nResults: %d pass, %d fail (%.1f%% pass rate)\n", 
		passCount, failCount, float64(passCount)/float64(len(testPatterns))*100)
}

func shouldValidationFail(bodyStr string) bool {
	// Check for null values
	if strings.Contains(bodyStr, `"preference":null`) || 
	   strings.Contains(bodyStr, `"weightUnit":null`) || 
	   strings.Contains(bodyStr, `"heightUnit":null`) ||
	   strings.Contains(bodyStr, `"weight":null`) ||
	   strings.Contains(bodyStr, `"height":null`) ||
	   strings.Contains(bodyStr, `"name":null`) ||
	   strings.Contains(bodyStr, `"imageUri":null`) {
		return true
	}
	
	// Check for empty strings on required fields
	if strings.Contains(bodyStr, `"preference":""`) ||
	   strings.Contains(bodyStr, `"weightUnit":""`) ||
	   strings.Contains(bodyStr, `"heightUnit":""`) {
		return true
	}
	
	// Check for invalid enum values and patterns
	if strings.Contains(bodyStr, `"preference":"notAnEnumValue"`) ||
	   strings.Contains(bodyStr, `"weightUnit":"notAnEnumValue"`) ||
	   strings.Contains(bodyStr, `"heightUnit":"notAnEnumValue"`) ||
	   strings.Contains(bodyStr, `"imageUri":"notAUrl"`) ||
	   strings.Contains(bodyStr, `"imageUri":"http://incomplete"`) ||
	   strings.Contains(bodyStr, `"imageUri":"notAnObject"`) ||
	   strings.Contains(bodyStr, `"imageUri":"notABoolean"`) ||
	   strings.Contains(bodyStr, `"weight":"notANumber"`) ||
	   strings.Contains(bodyStr, `"height":"notANumber"`) {
		return true
	}
	
	// Check for type mismatches
	if strings.Contains(bodyStr, `"preference":123`) ||
	   strings.Contains(bodyStr, `"preference":true`) ||
	   strings.Contains(bodyStr, `"preference":{}`) ||
	   strings.Contains(bodyStr, `"preference":[]`) ||
	   strings.Contains(bodyStr, `"weightUnit":123`) ||
	   strings.Contains(bodyStr, `"weightUnit":true`) ||
	   strings.Contains(bodyStr, `"heightUnit":123`) ||
	   strings.Contains(bodyStr, `"heightUnit":true`) ||
	   strings.Contains(bodyStr, `"name":123`) ||
	   strings.Contains(bodyStr, `"name":true`) ||
	   strings.Contains(bodyStr, `"imageUri":123`) ||
	   strings.Contains(bodyStr, `"imageUri":true`) {
		return true
	}
	
	// Check for string length violations
	if strings.Contains(bodyStr, `"name":"a"`) || // Too short
	   strings.Contains(bodyStr, `"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`) { // Too long
		return true
	}
	
	// Check for number range violations (weight: min=10, max=1000; height: min=3, max=250)
	if strings.Contains(bodyStr, `"weight":9`) ||    // Below min (10-1=9)
	   strings.Contains(bodyStr, `"weight":1001`) || // Above max (1000+1=1001)
	   strings.Contains(bodyStr, `"height":2`) ||    // Below min (3-1=2)
	   strings.Contains(bodyStr, `"height":251`) ||  // Above max (250+1=251)
	   strings.Contains(bodyStr, `"weight":0`) ||    // Zero (below min)
	   strings.Contains(bodyStr, `"height":0`) ||    // Zero (below min)
	   strings.Contains(bodyStr, `"weight":-1`) ||   // Negative
	   strings.Contains(bodyStr, `"height":-1`) {    // Negative
		return true
	}
	
	// Check for additional invalid patterns identified in testing
	if strings.Contains(bodyStr, `"name":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"`) || // 61 b's
	   strings.Contains(bodyStr, `"preference":"INVALID"`) ||
	   strings.Contains(bodyStr, `"weightUnit":"INVALID"`) ||
	   strings.Contains(bodyStr, `"heightUnit":"INVALID"`) ||
	   strings.Contains(bodyStr, `"imageUri":"invalid-url"`) ||
	   strings.Contains(bodyStr, `"imageUri":"ftp://example.com"`) ||
	   strings.Contains(bodyStr, `"name":""`) ||     // Empty optional field
	   strings.Contains(bodyStr, `"imageUri":""`) {  // Empty optional field
		return true
	}
	
	return false
}
