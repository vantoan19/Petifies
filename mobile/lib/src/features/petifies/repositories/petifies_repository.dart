import 'dart:io';

import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/petifies.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/petifies_service.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/address.dart' as addressModel;
import 'package:mobile/src/models/petifies.dart' as petifiesModel;
import 'package:mobile/src/proto/common/common.pb.dart' as commonProto;

final petifiesRepositoryProvider = Provider((ref) =>
    PetifiesRepository(petifiesService: ref.read(petifiesServiceProvider)));

abstract class IPetifiesRepository {
  Future<List<PetifiesModel>> listNearByPetifies({
    required petifiesModel.PetifiesType type,
    required double longitude,
    required double latitude,
    required double radius,
    int pageSize = 40,
    int offset = 0,
  });
  Future<PetifiesModel> createPetifies({
    required petifiesModel.PetifiesType type,
    required String title,
    required String description,
    required String petName,
    required List<NetworkImageModel> images,
    required addressModel.Address address,
  });
  Future<List<PetifiesModel>> listByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  });
}

class PetifiesRepository implements IPetifiesRepository {
  final PetifiesService _petifiesService;

  PetifiesRepository({required PetifiesService petifiesService})
      : this._petifiesService = petifiesService;

  petifiesModel.PetifiesType convertPetifiesProtoType(
      commonProto.PetifiesType type) {
    switch (type) {
      case commonProto.PetifiesType.PETIFIES_TYPE_CAT_ADOPTION:
        return petifiesModel.PetifiesType.CAT_ADOPTION;
      case commonProto.PetifiesType.PETIFIES_TYPE_CAT_PLAYING:
        return petifiesModel.PetifiesType.CAT_PLAYING;
      case commonProto.PetifiesType.PETIFIES_TYPE_CAT_SITTING:
        return petifiesModel.PetifiesType.CAT_SITTING;
      case commonProto.PetifiesType.PETIFIES_TYPE_DOG_ADOPTION:
        return petifiesModel.PetifiesType.DOG_ADOPTION;
      case commonProto.PetifiesType.PETIFIES_TYPE_DOG_SITTING:
        return petifiesModel.PetifiesType.DOG_SITTING;
      case commonProto.PetifiesType.PETIFIES_TYPE_DOG_WALKING:
        return petifiesModel.PetifiesType.DOG_WALKING;
      default:
        return petifiesModel.PetifiesType.UNKNOWN;
    }
  }

  commonProto.PetifiesType convertPetifiesType(
      petifiesModel.PetifiesType type) {
    switch (type) {
      case petifiesModel.PetifiesType.CAT_ADOPTION:
        return commonProto.PetifiesType.PETIFIES_TYPE_CAT_ADOPTION;
      case petifiesModel.PetifiesType.CAT_PLAYING:
        return commonProto.PetifiesType.PETIFIES_TYPE_CAT_PLAYING;
      case petifiesModel.PetifiesType.CAT_SITTING:
        return commonProto.PetifiesType.PETIFIES_TYPE_CAT_SITTING;
      case petifiesModel.PetifiesType.DOG_ADOPTION:
        return commonProto.PetifiesType.PETIFIES_TYPE_DOG_ADOPTION;
      case petifiesModel.PetifiesType.DOG_SITTING:
        return commonProto.PetifiesType.PETIFIES_TYPE_DOG_SITTING;
      case petifiesModel.PetifiesType.DOG_WALKING:
        return commonProto.PetifiesType.PETIFIES_TYPE_DOG_WALKING;
      default:
        return commonProto.PetifiesType.PETIFIES_TYPE_UNKNOWN;
    }
  }

  petifiesModel.PetifiesStatus convertProtoStatus(
      commonProto.PetifiesStatus status) {
    switch (status) {
      case commonProto.PetifiesStatus.PETIFIES_STATUS_DELETED:
        return petifiesModel.PetifiesStatus.DELETED;
      case commonProto.PetifiesStatus.PETIFIES_STATUS_AVAILABLE:
        return petifiesModel.PetifiesStatus.AVAILABLE;
      case commonProto.PetifiesStatus.PETIFIES_STATUS_UNAVAILABLE:
        return petifiesModel.PetifiesStatus.UNAVAILABLE;
      default:
        return petifiesModel.PetifiesStatus.UNKNOWN;
    }
  }

  PetifiesModel toPetifiesModel(PetifiesWithUserInfo petifies) {
    return PetifiesModel(
      id: petifies.id,
      owner: BasicUserInfoModel(
        id: petifies.owner.id,
        email: petifies.owner.email,
        firstName: petifies.owner.firstName,
        lastName: petifies.owner.lastName,
      ),
      type: convertPetifiesProtoType(petifies.type),
      title: petifies.title,
      description: petifies.description,
      petName: petifies.petName,
      image: petifies.images
          .map((e) => NetworkImageModel(uri: e.uri, description: e.description))
          .toList(),
      status: convertProtoStatus(petifies.status),
      address: addressModel.Address(
        addressLineOne: petifies.address.addressLineOne,
        addressLineTwo: petifies.address.addressLineTwo,
        street: petifies.address.street,
        district: petifies.address.district,
        city: petifies.address.city,
        region: petifies.address.region,
        postalCode: petifies.address.postalCode,
        country: petifies.address.country,
        longitude: petifies.address.longitude,
        latitude: petifies.address.latitude,
      ),
      createdAt: petifies.createdAt.toDateTime(),
    );
  }

  Future<List<PetifiesModel>> listNearByPetifies({
    required petifiesModel.PetifiesType type,
    required double longitude,
    required double latitude,
    required double radius,
    int pageSize = 40,
    int offset = 0,
  }) async {
    final petifies = await _petifiesService.listNearByPetifies(
      type: convertPetifiesType(type),
      longitude: longitude,
      latitude: latitude,
      radius: radius,
      pageSize: pageSize,
      offset: offset,
    );

    return petifies.petifies.map((e) => toPetifiesModel(e)).toList();
  }

  Future<PetifiesModel> createPetifies({
    required petifiesModel.PetifiesType type,
    required String title,
    required String description,
    required String petName,
    required List<NetworkImageModel> images,
    required addressModel.Address address,
  }) async {
    final petify = await _petifiesService.userCreatePetifies(
        type: convertPetifiesType(type),
        title: title,
        description: description,
        petName: petName,
        images: images,
        address: address);

    return toPetifiesModel(petify);
  }

  Future<List<PetifiesModel>> listByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    final petifies = await _petifiesService.listPetifiesByUserId(
      userId: userId,
      pageSize: pageSize,
      afterId: afterId,
    );

    return petifies.petifies.map((e) => toPetifiesModel(e)).toList();
  }
}
