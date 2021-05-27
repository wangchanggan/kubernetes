/*在Visitor Example代码示例中，定义了Visitor 接口，增加了VisitorList 对象，该对象相当于多个Visitor匿名函数的集合。
另外，增加了3个Visitor的类，分别实现Visit方法，在每一个VisitorFunc执行之前(before) 和执行之后(after) 分别输出print信息。
Visitor Example代码执行结果输出如下:
In Visitorl before fn
In VisitorList before fn
In Visitor2 before fn
In Visitor3 before fn
In visitFunc
In Visitor3 after fn
In Visitor2 after fn
In VisitorList after fn
In Visitorl after fn
通过Visitor代码示例的输出，能够更好地理解Visitor的多层嵌套关系。
根据输出结果，最先执行的是Visitor1 中fn 匿名函数之前的代码，然后是VisitorList、Visitor2 和Visitor3中 fn匿名函数之前的代码。
紧接着执行VisitFunc(visitor.Visit)。最后执行Visitor3、Visitor2、 VisitorList、 Visitorl 的fn匿名函数之后的代码。
整个多层嵌套关系的执行过程有些类似于递归操作。
*/

package main

import "fmt"

type Visitor interface {
	Visit(VisitorFunc) error
}

type VisitorFunc func() error

type VisitorList []Visitor

func (l VisitorList) Visitor(fn VisitorFunc) error {
	for i := range l {
		if err := l[i].Visit(func() error {
			fmt.Println("In VisitorList before fn")
			fn()
			fmt.Println("In VisitorList after fn")
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

type Visitorl struct {
}

func (v Visitorl) Visit(fn VisitorFunc) error {
	fmt.Println("In Visitorl before fn")
	fn()
	fmt.Println("In Visitorl after fn")
	return nil
}

type Visitor2 struct {
	visitorList VisitorList
}

func (v Visitor2) Visit(fn VisitorFunc) error {
	v.visitorList.Visitor(func() error {
		fmt.Println("In Visitor2 before fn")
		fn()
		fmt.Println("In Visitor2 after fn")
		return nil
	})
	return nil
}

type Visitor3 struct {
	visitor Visitor
}

func (v Visitor3) Visit(fn VisitorFunc) error {
	v.visitor.Visit(func() error {
		fmt.Println("In Visitor3 before fn")
		fn()
		fmt.Println("In Visitor3 after fn")
		return nil
	})
	return nil
}

func main() {
	var visitor Visitor
	var visitors []Visitor

	/*在main函数中，首先将Visitor1嵌入VisitorList中，VisitorList 是Visitor的集合，可存放多个Visitor。
	然后将VisitorList 嵌入Visitor2中，接着将Visitor2嵌入Visitor3 中。
	最终形成Visitor3 {Visitor2 {VisitorList{Visitorl}}}的嵌套关系。*/
	visitor = Visitorl{}
	visitors = append(visitors, visitor)
	visitor = Visitor2{VisitorList(visitors)}
	visitor = Visitor3{visitor}
	visitor.Visit(func() error {
		fmt.Println("In visitFunc")
		return nil
	})
}
