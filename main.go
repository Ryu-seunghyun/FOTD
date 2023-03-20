// ( 인터페이스, 가변 인자, defer, closure, map
//  에러핸들링, 고루틴( sync , 채널,  select ),
//  컨텍스트는 제외

/* 점.메.추 기능 구현
   B1, 1, Out  (대)
   C, A, K, J   (중)
   rice, noodle, fastfood, stew (소)
*/

/* 사람 - 인터페이스   {파트별}
  	ㄴ (메서드) 수동 선택
	ㄴ (메서드) 자동 선택
*/

/* 수동 선택(메서드) - 고루틴
(채널) 대분류 정하기  (대_채-공급,해제)
(채널) 중분류 정하기  (대_채-소비,중_채-공급,해제)
(채널) 대분류 정하기  (중_채-소비)
>> func Customize(대,중,소)  < < 각각 입력받아서 루틴 처리
*/
/* 자동 선택(메서드) - 클로저
분류별 랜덤 생성 -> 수동 선택 인자로 할당 및 호출
*/

// Err_Code 1 : 비정상 입력 패닉띄운 후 defer 처리로 값 재입력 받고 재실행
// Err_Code 2 : 비정상 동작 오류, 프로그램 종료
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"golang.org/x/exp/slices"
)

// 임시 데이터
var (
	Places       = []string{"B1", "1", "Out"}
	Styles       = []string{"C", "A", "K", "J"}
	F_Categories = []string{"rice", "noodle", "fastfood", "strew"}
)

type FOTD struct { // Food of the day  (Food Model)
	// ID
	Place      string
	Style      string
	F_Category string
}

type Customer interface {
	Recommand() *FOTD            // <- 자동 선택
	Customize(interface{}) *FOTD // <- 수동 선택
}

type Backend struct { // 백엔드팀 객체,  ... 프론트팀, 기획팀 ...
	Name    string
	Emp_num string
	Menu    *FOTD
	// Details interface{}
}

func RandomPlay(c Customer) {
	c.Recommand()
}
func CustomPlay(c Customer, order interface{}) {
	c.Customize(order)
}

// wg, ch , select ? ,
func (b *Backend) Recommand() *FOTD {
	rd_num := rand.Intn(899) + 100
	fotd := b.Customize(rd_num)
	return fotd
}
func (b *Backend) Customize(order interface{}) *FOTD { // int, []string
	P_ch := make(chan interface{})
	S_ch := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(3)
	go b.Menu.SelectPlace(order, P_ch, &wg)
	go b.Menu.SelectStyle(P_ch, S_ch, &wg)
	go b.Menu.SelectCategory(S_ch, &wg)
	// close(P_ch)
	// close(S_ch)
	wg.Wait()
	return b.Menu
}

func (f *FOTD) SelectPlace(order interface{}, P_ch chan interface{}, wg *sync.WaitGroup) error {
	// 숫자로 입력받거나 문자열 여러개로 입력 ( )  (문자열 여러개는 문자열 슬라이스로 캐스팅)
	// 3자리 이하의 숫자 및 문자열 슬라이스만 입력 가능 / (입력받지 못한 부분은 임의로 추천 ,문자열 중 선택지 목록에 없는 경우 추천_(미구현)))
	// 0_0_0   [Place_Style_F-category]

	switch order.(type) {
	case int: // 3자리 수 확인
		if int(order.(int)/1000) > 0 {
			// Err_Code 1
			return errors.New("occured error 1")
		}
		var Place_num = int(order.(int) / 100) // 첫 번째 숫자 (0)
		// log.Panicf("RandomPick : %v", Places[RandomPick(Places)])
		// log.Panicf("Place_num : %v", Place_num)
		// log.Panicf("Place Num : %v", Places[Place_num%len(Places)])
		// log.Panicf("Order info : %v", order)

		var Mod_num = int(order.(int) % 100) // 첫 번째 숫자를 제외한 나머지 숫자 (00)

		if Place_num == 0 || Place_num > len(Places) { // 0 이거나 len(Places) 이상 경우 선택장애로 판단,  추천
			f.Place = Places[RandomPick(Places)] // nil pointer | invalid memory addr . 빈값에 접근
			P_ch <- Mod_num
			wg.Done()
			return nil
		}
		f.Place = Places[Place_num%len(Places)]

		P_ch <- Mod_num
		wg.Done()
		return nil

	case []string:
		if !slices.Contains(Places, order.([]string)[0]) {
			f.Place = Places[RandomPick(Places)]
		} else {
			f.Place = order.([]string)[0]
		}
		P_ch <- order.([]string)[1:]
		wg.Done()
		return nil
	case interface{}:
		// Err_Code 1
		wg.Done()
		return errors.New("occured error 2")
	}
	wg.Done()
	return errors.New("occured error 3")
}

func RandomPick(i []string) int {
	return int(rand.Intn(len(i) - 1))
}

func (f *FOTD) SelectStyle(P_ch, S_ch chan interface{}, wg *sync.WaitGroup) error {
	order := <-P_ch
	switch order.(type) {
	case int:
		var Style_num = int(order.(int) / 10) // 두 번째 숫자 (0)
		var Mod_num = int(order.(int) % 10)   // 두 번째 숫자를 제외한 나머지 숫자 (0)

		if Style_num == 0 || Style_num > len(Styles) { // 0 이거나 len(Styles) 이상 경우 선택장애로 판단,  추천
			f.Style = Styles[RandomPick(Styles)]
		} else {
			f.Style = Styles[int(order.(int)/100)%len(Styles)]
		}
		S_ch <- Mod_num
		wg.Done()
		return nil

	case []string:
		if !slices.Contains(Styles, order.([]string)[0]) {
			f.Style = Styles[RandomPick(Styles)]
		} else {
			f.Style = order.([]string)[0]
		}
		S_ch <- order.([]string)[1:]
		wg.Done()
		return nil
	default:
		// Err_Code 2
		return errors.New("occured error 4")
	}
}

func (f *FOTD) SelectCategory(S_ch chan interface{}, wg *sync.WaitGroup) error {
	order := <-S_ch
	switch order.(type) {
	case int:
		Category_num := order.(int)
		if Category_num == 0 || Category_num > len(F_Categories) {
			f.F_Category = F_Categories[RandomPick(F_Categories)]
		} else {
			f.F_Category = F_Categories[Category_num%len(F_Categories)]
		}
		wg.Done()
		return nil

	case []string:
		if !slices.Contains(F_Categories, order.([]string)[0]) {
			f.F_Category = F_Categories[RandomPick(F_Categories)]
		} else {
			f.F_Category = order.([]string)[0]
		}
		wg.Done()
		return nil
	default:
		// Err_Code 2
		return errors.New("occured error 5")
	}
}

func main() {
	ryu := &Backend{
		Name:    "RYU",
		Emp_num: "12303510",
		Menu:    &FOTD{},
	}
	RandomPlay(ryu)
	fmt.Println(*ryu.Menu)
}
