type A struct{
    v: u8
}

func (A) getV()u8{
    return self.v
}

func test()A{
    return {11}
}

@extern(main)
func main()u8{
    return test().getV()
}
