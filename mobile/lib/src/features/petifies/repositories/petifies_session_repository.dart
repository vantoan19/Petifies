import 'dart:io';

import 'package:mobile/src/models/petifies_session.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/petifies_service.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/petifies_session.dart'
    as petifiesSessionModel;
import 'package:mobile/src/proto/common/common.pb.dart' as commonProto;

final petifiesSessionRepositoryProvider = Provider((ref) =>
    PetifiesSessionRepository(
        petifiesService: ref.read(petifiesServiceProvider)));

abstract class IPetifiesSessionRepository {
  Future<PetifiesSessionModel> createSession({
    required String petifiesId,
    required DateTime fromTime,
    required DateTime toTime,
  });
  Future<List<PetifiesSessionModel>> listSessionsByPetifiesId({
    required String petifiesId,
    int pageSize = 40,
    String afterId = "",
  });
}

class PetifiesSessionRepository implements IPetifiesSessionRepository {
  final PetifiesService _petifiesService;

  PetifiesSessionRepository({required PetifiesService petifiesService})
      : this._petifiesService = petifiesService;

  petifiesSessionModel.PetifiesSessionStatus convertProtoStatus(
      commonProto.PetifiesSessionStatus status) {
    switch (status) {
      case commonProto
          .PetifiesSessionStatus.PETIFIES_SESSION_STATUS_WAITING_FOR_PROPOSAL:
        return petifiesSessionModel.PetifiesSessionStatus.WAITING_FOR_PROPOSAL;
      case commonProto
          .PetifiesSessionStatus.PETIFIES_SESSION_STATUS_PROPOSAL_ACCEPTED:
        return petifiesSessionModel.PetifiesSessionStatus.PROPOSAL_ACCEPTED;
      case commonProto.PetifiesSessionStatus.PETIFIES_SESSION_STATUS_ON_GOING:
        return petifiesSessionModel.PetifiesSessionStatus.ON_GOING;
      case commonProto.PetifiesSessionStatus.PETIFIES_SESSION_STATUS_ENDED:
        return petifiesSessionModel.PetifiesSessionStatus.ENDED;
      default:
        return petifiesSessionModel.PetifiesSessionStatus.UNKNOWN;
    }
  }

  PetifiesSessionModel toPetifiesSessionModel(UserPetifiesSession session) {
    return PetifiesSessionModel(
      id: session.id,
      petifiesId: session.petifiesId,
      fromTime: session.fromTime.toDateTime(),
      toTime: session.toTime.toDateTime(),
      status: convertProtoStatus(session.status),
      createdAt: session.createdAt.toDateTime(),
    );
  }

  Future<PetifiesSessionModel> createSession({
    required String petifiesId,
    required DateTime fromTime,
    required DateTime toTime,
  }) async {
    final session = await _petifiesService.userCreatePetifiesSession(
      petifiesId: petifiesId,
      fromTime: fromTime,
      toTime: toTime,
    );

    return toPetifiesSessionModel(session);
  }

  Future<List<PetifiesSessionModel>> listSessionsByPetifiesId({
    required String petifiesId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    final sessions = await _petifiesService.listSessionsByPetifiesId(
      petifiesId: petifiesId,
      pageSize: pageSize,
      afterId: afterId,
    );

    return sessions.sessions.map((e) => toPetifiesSessionModel(e)).toList();
  }
}
