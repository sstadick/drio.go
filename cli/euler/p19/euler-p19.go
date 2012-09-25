package main

import "fmt"
import "time"

// Problem 19
//
//    14 June 2002
//
//    You are given the following information, but you may prefer to do some
//    research for yourself.
//      * 1 Jan 1900 was a Monday.
//      * Thirty days has September,
//        April, June and November.
//        All the rest have thirty-one,
//        Saving February alone,
//        Which has twenty-eight, rain or shine.
//        And on leap years, twenty-nine.
//      * A leap year occurs on any year evenly divisible by 4, but not on a
//        century unless it is divisible by 400.
//
//    How many Sundays fell on the first of the month during the twentieth
//    century (1 Jan 1901 to 31 Dec 2000)?
func main() {
  pMonth := time.January
  sunday := time.Sunday
  december := time.December
  t := time.Date(1901, 1, 0, 0, 0, 0, 0, time.UTC)
  sundays := 0
  for {
    if t.Month() != pMonth && t.Weekday() == sunday {
      sundays++
    }
    if t.Year() == 2000 && t.Month() == december && t.Day() == 31 {
      break
    }
    pMonth = t.Month()
    t = t.AddDate(0, 0, 1)
  }
  fmt.Println(sundays)
}
