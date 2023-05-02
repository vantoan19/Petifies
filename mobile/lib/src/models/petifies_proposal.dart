// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'package:mobile/src/models/basic_user_info.dart';

enum PetifiesProposalStatus {
  WAITING_FOR_ACCEPTANCE,
  ACCEPTED,
  CANCELLED,
  REJECTED,
  SESSION_CLOSED,
  UNKNOWN
}

class PetifiesProposalModel {
  final String id;
  final BasicUserInfoModel author;
  final String petifiesSessionId;
  final String proposal;
  final PetifiesProposalStatus status;
  final DateTime createdAt;
  PetifiesProposalModel({
    required this.id,
    required this.author,
    required this.petifiesSessionId,
    required this.proposal,
    required this.status,
    required this.createdAt,
  });

  PetifiesProposalModel copyWith({
    String? id,
    BasicUserInfoModel? author,
    String? petifiesSessionId,
    String? proposal,
    PetifiesProposalStatus? status,
    DateTime? createdAt,
  }) {
    return PetifiesProposalModel(
      id: id ?? this.id,
      author: author ?? this.author,
      petifiesSessionId: petifiesSessionId ?? this.petifiesSessionId,
      proposal: proposal ?? this.proposal,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'PetifiesProposalModel(id: $id, author: $author, petifiesSessionId: $petifiesSessionId, proposal: $proposal, status: $status, createdAt: $createdAt)';
  }

  @override
  bool operator ==(covariant PetifiesProposalModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.author == author &&
        other.petifiesSessionId == petifiesSessionId &&
        other.proposal == proposal &&
        other.status == status &&
        other.createdAt == createdAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        author.hashCode ^
        petifiesSessionId.hashCode ^
        proposal.hashCode ^
        status.hashCode ^
        createdAt.hashCode;
  }
}
