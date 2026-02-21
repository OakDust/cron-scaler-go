package converter

import (
	"scale-handler/internal/domain"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func DomainToProto(schedule *domain.Schedule) *scalehandlerv1.Schedule {
	if schedule == nil {
		return nil
	}

	protoSchedule := &scalehandlerv1.Schedule{
		Weekdays:   make(map[string]*scalehandlerv1.Schedule_DaySchedule),
		Dates:      make(map[string]*scalehandlerv1.Schedule_DaySchedule),
		Exceptions: schedule.Rules.Exceptions,
	}

	// Конвертируем weekdays
	for day, ranges := range schedule.Rules.Weekdays {
		daySchedule := &scalehandlerv1.Schedule_DaySchedule{}
		for _, tr := range ranges {
			daySchedule.TimeRanges = append(daySchedule.TimeRanges, &scalehandlerv1.TimeRange{
				From:     tr.From,
				To:       tr.To,
				Replicas: tr.Replicas,
			})
		}
		protoSchedule.Weekdays[day] = daySchedule
	}

	// Конвертируем dates
	for date, ranges := range schedule.Rules.Dates {
		daySchedule := &scalehandlerv1.Schedule_DaySchedule{}
		for _, tr := range ranges {
			daySchedule.TimeRanges = append(daySchedule.TimeRanges, &scalehandlerv1.TimeRange{
				From:     tr.From,
				To:       tr.To,
				Replicas: tr.Replicas,
			})
		}
		protoSchedule.Dates[date] = daySchedule
	}

	return protoSchedule
}

func ProtoToDomainRules(protoSchedule *scalehandlerv1.Schedule) domain.ScheduleRules {
	if protoSchedule == nil {
		return domain.ScheduleRules{}
	}

	rules := domain.ScheduleRules{
		Weekdays:   make(map[string][]domain.TimeRange),
		Dates:      make(map[string][]domain.TimeRange),
		Exceptions: protoSchedule.Exceptions,
	}

	// Конвертируем weekdays
	for day, daySchedule := range protoSchedule.Weekdays {
		var ranges []domain.TimeRange
		if daySchedule != nil {
			for _, tr := range daySchedule.TimeRanges {
				ranges = append(ranges, domain.TimeRange{
					From:     tr.From,
					To:       tr.To,
					Replicas: tr.Replicas,
				})
			}
		}
		rules.Weekdays[day] = ranges
	}

	// Конвертируем dates
	for date, daySchedule := range protoSchedule.Dates {
		var ranges []domain.TimeRange
		if daySchedule != nil {
			for _, tr := range daySchedule.TimeRanges {
				ranges = append(ranges, domain.TimeRange{
					From:     tr.From,
					To:       tr.To,
					Replicas: tr.Replicas,
				})
			}
		}
		rules.Dates[date] = ranges
	}

	return rules
}
