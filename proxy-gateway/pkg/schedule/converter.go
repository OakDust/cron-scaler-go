package schedule

import (
	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
)

// DTOToProto конвертирует REST DTO в gRPC protobuf
func DTOToProto(dto *ScheduleDTO) *scalehandlerv1.Schedule {
	if dto == nil {
		return nil
	}

	proto := &scalehandlerv1.Schedule{
		Weekdays:   make(map[string]*scalehandlerv1.Schedule_DaySchedule),
		Dates:      make(map[string]*scalehandlerv1.Schedule_DaySchedule),
		Exceptions: dto.Exceptions,
	}

	for day, ranges := range dto.Weekdays {
		proto.Weekdays[day] = &scalehandlerv1.Schedule_DaySchedule{
			TimeRanges: timeRangesToProto(ranges),
		}
	}

	for date, ranges := range dto.Dates {
		proto.Dates[date] = &scalehandlerv1.Schedule_DaySchedule{
			TimeRanges: timeRangesToProto(ranges),
		}
	}

	return proto
}

// ProtoToDTO конвертирует gRPC protobuf в REST DTO
func ProtoToDTO(proto *scalehandlerv1.Schedule) *ScheduleDTO {
	if proto == nil {
		return nil
	}

	dto := &ScheduleDTO{
		Weekdays:   make(map[string][]TimeRangeDTO),
		Dates:      make(map[string][]TimeRangeDTO),
		Exceptions: proto.Exceptions,
	}

	for day, daySchedule := range proto.Weekdays {
		if daySchedule != nil {
			dto.Weekdays[day] = timeRangesToDTO(daySchedule.TimeRanges)
		} else {
			dto.Weekdays[day] = []TimeRangeDTO{}
		}
	}

	for date, daySchedule := range proto.Dates {
		if daySchedule != nil {
			dto.Dates[date] = timeRangesToDTO(daySchedule.TimeRanges)
		} else {
			dto.Dates[date] = []TimeRangeDTO{}
		}
	}

	return dto
}

func timeRangesToProto(ranges []TimeRangeDTO) []*scalehandlerv1.TimeRange {
	if ranges == nil {
		return nil
	}
	result := make([]*scalehandlerv1.TimeRange, len(ranges))
	for i, r := range ranges {
		result[i] = &scalehandlerv1.TimeRange{
			From:     r.From,
			To:       r.To,
			Replicas: r.Replicas,
		}
	}
	return result
}

func timeRangesToDTO(ranges []*scalehandlerv1.TimeRange) []TimeRangeDTO {
	if ranges == nil {
		return nil
	}
	result := make([]TimeRangeDTO, len(ranges))
	for i, r := range ranges {
		if r != nil {
			result[i] = TimeRangeDTO{From: r.From, To: r.To, Replicas: r.Replicas}
		}
	}
	return result
}
