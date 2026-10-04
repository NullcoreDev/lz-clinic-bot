package catalog

type Service struct {
	ID       int64
	Name     string
	Duration int
}

type Doctor struct {
	ID             int64
	Name           string
	Specialization string
	ServiceIDs     []int64
}

var services = []Service{
	{
		ID:       1,
		Name:     "Приём врача-невролога",
		Duration: 60,
	},
	{
		ID:       2,
		Name:     "Приём врача-кардиолога",
		Duration: 60,
	},
	{
		ID:       3,
		Name:     "Приём врача-терапевта",
		Duration: 60,
	},
}

var doctors = []Doctor{
	{
		ID:             1,
		Name:           "Раков Александр Михайлович",
		Specialization: "Врач-невролог",
		ServiceIDs:     []int64{1},
	},
	{
		ID:             2,
		Name:           "Гуревич Оксана Васильевна",
		Specialization: "Врач-кардиолог",
		ServiceIDs:     []int64{2},
	},
	{
		ID:             3,
		Name:           "Игнатенкова Эльвира Ильгизовна",
		Specialization: "Врач-терапевт",
		ServiceIDs:     []int64{3},
	},
}

func GetServices() []Service {
	return services
}

func GetDoctors() []Doctor {
	return doctors
}

func GetServiceByID(id int64) (Service, bool) {
	for _, service := range services {
		if service.ID == id {
			return service, true
		}
	}

	return Service{}, false
}

func GetDoctorByID(id int64) (Doctor, bool) {
	for _, doctor := range doctors {
		if doctor.ID == id {
			return doctor, true
		}
	}

	return Doctor{}, false
}

func GetDoctorsByService(serviceID int64) []Doctor {
	var result []Doctor

	for _, doctor := range doctors {
		for _, id := range doctor.ServiceIDs {
			if id == serviceID {
				result = append(result, doctor)
				break
			}
		}
	}

	return result
}
