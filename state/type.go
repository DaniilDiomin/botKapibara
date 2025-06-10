package state

type State int

const (
	Idle State = iota // Начальное состояние

	RestarauntSelection // Выбор ресторана

	FreshcoffSelection //ресторан ФРЕШКОФФ
	RogachevSelection  //ресторан Рогачев
	RechicaSelection   //ресторан речица

	RequestSubmission //Оформление заявки
	WorkSchedule      //Учет рабочего времени
	WriteOff          //Списание товаров

	SurveyInProgress
	RoleSelection
)

type stateStruct struct {
	Idle struct {
		RestarauntSelection struct {
			FreshcoffSelection struct {
			}
			RogachevSelection struct {
			}
			RechicaSelection struct {
			}
		}
	}
}

type UserState struct {
	Current State
	Context map[string]string
}
type StateManager struct {
	states map[int64]*UserState
}
