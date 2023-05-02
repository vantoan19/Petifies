import 'package:mobile/src/features/petifies/repositories/petifies_proposal_repository.dart';
import 'package:mobile/src/features/petifies/repositories/petifies_session_repository.dart';
import 'package:mobile/src/models/petifies_proposal.dart';
import 'package:mobile/src/models/petifies_session.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'petifies_session_controllers.g.dart';

@riverpod
class CreatePetifiesSessionController
    extends _$CreatePetifiesSessionController {
  @override
  void build() {}

  Future<PetifiesSessionModel> createPetifiesSession({
    required String petifiesId,
    required DateTime fromTime,
    required DateTime toTime,
  }) async {
    final sessionRepository = ref.read(petifiesSessionRepositoryProvider);
    return await sessionRepository.createSession(
      petifiesId: petifiesId,
      fromTime: fromTime,
      toTime: toTime,
    );
  }
}

@riverpod
class ListSessionsByPetifiesIdController
    extends _$ListSessionsByPetifiesIdController {
  Future<List<PetifiesSessionModel>> build({required String petifiesId}) async {
    final sessionRepository = ref.read(petifiesSessionRepositoryProvider);

    return await sessionRepository.listSessionsByPetifiesId(
        petifiesId: petifiesId);
  }

  void fetchMoreSessions() {
    state.whenData((value) async {
      state = AsyncLoading();
      state = await AsyncValue.guard(() async {
        final sessionRepository = ref.read(petifiesSessionRepositoryProvider);
        final sessions = await sessionRepository.listSessionsByPetifiesId(
            petifiesId: petifiesId, afterId: value.last.id);
        return [...value, ...sessions];
      });
    });
  }

  void addSession(PetifiesSessionModel session) {
    state.whenData((value) {
      state = AsyncData([...value, session]);
    });
  }
}
