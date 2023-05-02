// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

enum PetifiesSessionStatus {
  WAITING_FOR_PROPOSAL,
  PROPOSAL_ACCEPTED,
  ON_GOING,
  ENDED,
  UNKNOWN
}

class PetifiesSessionModel {
  final String id;
  final String petifiesId;
  final DateTime fromTime;
  final DateTime toTime;
  final PetifiesSessionStatus status;
  final DateTime createdAt;
  PetifiesSessionModel({
    required this.id,
    required this.petifiesId,
    required this.fromTime,
    required this.toTime,
    required this.status,
    required this.createdAt,
  });

  PetifiesSessionModel copyWith({
    String? id,
    String? petifiesId,
    DateTime? fromTime,
    DateTime? toTime,
    PetifiesSessionStatus? status,
    DateTime? createdAt,
  }) {
    return PetifiesSessionModel(
      id: id ?? this.id,
      petifiesId: petifiesId ?? this.petifiesId,
      fromTime: fromTime ?? this.fromTime,
      toTime: toTime ?? this.toTime,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'PetifiesSessionModel(id: $id, petifiesId: $petifiesId, fromTime: $fromTime, toTime: $toTime, status: $status, createdAt: $createdAt)';
  }

  @override
  bool operator ==(covariant PetifiesSessionModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.petifiesId == petifiesId &&
        other.fromTime == fromTime &&
        other.toTime == toTime &&
        other.status == status &&
        other.createdAt == createdAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        petifiesId.hashCode ^
        fromTime.hashCode ^
        toTime.hashCode ^
        status.hashCode ^
        createdAt.hashCode;
  }
}
