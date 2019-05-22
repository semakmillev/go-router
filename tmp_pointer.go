package main

import "fmt"

func innerChange(v *string){
	*v = *v + "_inner_changed"
}

func change(v string){
	innerChange(&v)
	fmt.Println(v)
}

func main()  {
	a := "1234"
	fmt.Println(a)
	change(a)
	fmt.Println("after all",a)
}
