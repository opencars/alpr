package service

type CustomerService struct{}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (svc *CustomerService) Recognize() {

}
