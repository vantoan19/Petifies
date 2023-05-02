import 'package:mobile/src/features/petifies/repositories/petifies_proposal_repository.dart';
import 'package:mobile/src/models/petifies_proposal.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'petifies_proposal_controllers.g.dart';

@riverpod
class CreatePetifiesProposalController
    extends _$CreatePetifiesProposalController {
  @override
  void build() {}

  Future<PetifiesProposalModel> createPetifiesProposal({
    required String petifiesSessionId,
    required String proposal,
  }) async {
    final proposalRepository = ref.read(petifiesProposalRepositoryProvider);
    return await proposalRepository.createProposal(
        petifiesSessionId: petifiesSessionId, proposal: proposal);
  }
}

@riverpod
class ListProposalsByUserIdController
    extends _$ListProposalsByUserIdController {
  @override
  Future<List<PetifiesProposalModel>> build({required String userId}) async {
    final proposalRepository = ref.read(petifiesProposalRepositoryProvider);
    return await proposalRepository.listProposalsByUserId(userId: userId);
  }

  void fetchMoreProposals() {
    state.whenData((value) async {
      state = AsyncLoading();
      state = await AsyncValue.guard(() async {
        final proposalRepository = ref.read(petifiesProposalRepositoryProvider);
        final proposals = await proposalRepository.listProposalsByUserId(
            userId: userId, afterId: value.last.id);
        return [...value, ...proposals];
      });
    });
  }
}

@riverpod
class ListProposalsBySessionIdController
    extends _$ListProposalsBySessionIdController {
  @override
  Future<List<PetifiesProposalModel>> build({required String sessionId}) async {
    final proposalRepository = ref.read(petifiesProposalRepositoryProvider);
    return await proposalRepository.listProposalsBySessionId(
      sessionId: sessionId,
    );
  }

  void fetchMoreProposals() {
    state.whenData((value) async {
      state = AsyncLoading();
      state = await AsyncValue.guard(() async {
        final proposalRepository = ref.read(petifiesProposalRepositoryProvider);
        final proposals = await proposalRepository.listProposalsBySessionId(
            sessionId: sessionId, afterId: value.last.id);
        return [...value, ...proposals];
      });
    });
  }
}
