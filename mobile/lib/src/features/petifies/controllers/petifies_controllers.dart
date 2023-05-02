import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/petifies/repositories/petifies_repository.dart';
import 'package:mobile/src/features/petifies/screens/petifies_explore_screen.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/petifies.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:mobile/src/models/address.dart' as addressModel;
import 'package:mobile/src/models/petifies.dart' as petifiesModel;
import 'package:freezed_annotation/freezed_annotation.dart';

part 'petifies_controllers.g.dart';
part 'petifies_controllers.freezed.dart';

@freezed
abstract class ListPetifiesParameters with _$ListPetifiesParameters {
  factory ListPetifiesParameters({
    required PetifiesType type,
    required double longitude,
    required double latitude,
    required double radius,
    required int pageSize,
    required int offset,
    required bool isMapConsumer,
  }) = _ListPetifiesParameters;
}

@riverpod
class ListPetifiesController extends _$ListPetifiesController {
  @override
  Future<List<PetifiesModel>> build(
      {required ListPetifiesParameters parameters}) async {
    final petifiesRepository = ref.read(petifiesRepositoryProvider);

    final petifies = await petifiesRepository.listNearByPetifies(
      type: parameters.type,
      longitude: parameters.longitude,
      latitude: parameters.latitude,
      radius: parameters.radius,
      pageSize: parameters.pageSize,
      offset: parameters.offset,
    );

    if (!parameters.isMapConsumer) {
      ref
          .read(petifiesListProvider(parameters.type, parameters.isMapConsumer)
              .notifier)
          .addPetifies(petifies);
    }
    return petifies;
  }

  void fetchNextPetifies(int offset) async {
    state = AsyncLoading();
    state = await AsyncValue.guard(() async {
      final petifiesRepository = ref.read(petifiesRepositoryProvider);
      final fetchedPetifies = await petifiesRepository.listNearByPetifies(
          type: parameters.type,
          longitude: parameters.longitude,
          latitude: parameters.latitude,
          radius: parameters.radius,
          offset: offset);
      if (!parameters.isMapConsumer) {
        ref
            .read(
                petifiesListProvider(parameters.type, parameters.isMapConsumer)
                    .notifier)
            .addPetifies(fetchedPetifies);
      }
      return fetchedPetifies;
    });
  }
}

@riverpod
class CreatePetifiesController extends _$CreatePetifiesController {
  @override
  void build() {}

  Future<PetifiesModel> createPetifies({
    required petifiesModel.PetifiesType type,
    required String title,
    required String description,
    required String petName,
    required List<Future<Either<Failure, NetworkImageModel>>> imageFutures,
    required addressModel.Address address,
  }) async {
    final imageEithers = await Future.wait(imageFutures);
    List<NetworkImageModel> images = [];
    imageEithers.forEach((either) {
      either.fold((l) => null, (r) => images.add(r));
    });

    final petifiesRepository = ref.read(petifiesRepositoryProvider);
    return await petifiesRepository.createPetifies(
        type: type,
        title: title,
        description: description,
        petName: petName,
        images: images,
        address: address);
  }
}

@riverpod
class ListPetifiesByUserId extends _$ListPetifiesByUserId {
  @override
  Future<List<PetifiesModel>> build({required String userId}) async {
    final petifiesRepository = ref.read(petifiesRepositoryProvider);

    return await petifiesRepository.listByUserId(userId: userId);
  }

  void fetchMorePetifies() {
    state.whenData((value) async {
      state = AsyncLoading();
      state = await AsyncValue.guard(() async {
        final petifiesRepository = ref.read(petifiesRepositoryProvider);
        final petifies = await petifiesRepository.listByUserId(
            userId: userId, afterId: value.last.id);
        return [...value, ...petifies];
      });
    });
  }
}
