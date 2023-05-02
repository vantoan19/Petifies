import 'dart:io';

import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/petifies_proposal.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/petifies_service.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/petifies_proposal.dart'
    as petifiesProposalModel;
import 'package:mobile/src/proto/common/common.pb.dart' as commonProto;

final petifiesProposalRepositoryProvider = Provider((ref) =>
    PetifiesProposalRepository(
        petifiesService: ref.read(petifiesServiceProvider)));

abstract class IPetifiesProposalRepository {
  Future<PetifiesProposalModel> createProposal({
    required String petifiesSessionId,
    required String proposal,
  });
  Future<List<petifiesProposalModel.PetifiesProposalModel>>
      listProposalsByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  });
  Future<List<petifiesProposalModel.PetifiesProposalModel>>
      listProposalsBySessionId({
    required String sessionId,
    int pageSize = 40,
    String afterId = "",
  });
}

class PetifiesProposalRepository implements IPetifiesProposalRepository {
  final PetifiesService _petifiesService;

  PetifiesProposalRepository({required PetifiesService petifiesService})
      : this._petifiesService = petifiesService;

  petifiesProposalModel.PetifiesProposalStatus convertProtoStatus(
      commonProto.PetifiesProposalStatus status) {
    switch (status) {
      case commonProto.PetifiesProposalStatus
          .PETIFIES_PROPOSAL_STATUS_WAITING_FOR_ACCEPTANCE:
        return petifiesProposalModel
            .PetifiesProposalStatus.WAITING_FOR_ACCEPTANCE;
      case commonProto.PetifiesProposalStatus.PETIFIES_PROPOSAL_STATUS_ACCEPTED:
        return petifiesProposalModel.PetifiesProposalStatus.ACCEPTED;
      case commonProto
          .PetifiesProposalStatus.PETIFIES_PROPOSAL_STATUS_CANCELLED:
        return petifiesProposalModel.PetifiesProposalStatus.CANCELLED;
      case commonProto
          .PetifiesProposalStatus.PETIFIES_PROPOSAL_STATUS_SESSION_CLOSED:
        return petifiesProposalModel.PetifiesProposalStatus.SESSION_CLOSED;
      case commonProto.PetifiesProposalStatus.PETIFIES_PROPOSAL_STATUS_REJECTED:
        return petifiesProposalModel.PetifiesProposalStatus.REJECTED;
      default:
        return petifiesProposalModel.PetifiesProposalStatus.UNKNOWN;
    }
  }

  PetifiesProposalModel toPetifiesProposalModel(
      PetifiesProposalWithUserInfo proposal) {
    return PetifiesProposalModel(
      id: proposal.id,
      author: BasicUserInfoModel(
        id: proposal.user.id,
        email: proposal.user.email,
        firstName: proposal.user.firstName,
        lastName: proposal.user.lastName,
      ),
      proposal: proposal.proposal,
      petifiesSessionId: proposal.petifiesSessionId,
      status: convertProtoStatus(proposal.status),
      createdAt: proposal.createdAt.toDateTime(),
    );
  }

  Future<PetifiesProposalModel> createProposal({
    required String petifiesSessionId,
    required String proposal,
  }) async {
    final proposalResp = await _petifiesService.userCreatePetifiesProposal(
      petifiesSessionId: petifiesSessionId,
      proposal: proposal,
    );

    return toPetifiesProposalModel(proposalResp);
  }

  Future<List<petifiesProposalModel.PetifiesProposalModel>>
      listProposalsByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    final proposals = await _petifiesService.listProposalsByUserId(
      userId: userId,
      pageSize: pageSize,
      afterId: afterId,
    );

    return proposals.proposals.map((e) => toPetifiesProposalModel(e)).toList();
  }

  Future<List<petifiesProposalModel.PetifiesProposalModel>>
      listProposalsBySessionId({
    required String sessionId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    final proposals = await _petifiesService.listProposalsBySessionId(
      sessionId: sessionId,
      pageSize: pageSize,
      afterId: afterId,
    );

    return proposals.proposals.map((e) => toPetifiesProposalModel(e)).toList();
  }
}
