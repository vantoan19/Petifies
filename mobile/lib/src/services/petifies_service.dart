import 'dart:async';

import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/address.dart' as addressModel;
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/proto/google/protobuf/timestamp.pb.dart';
import 'package:mobile/src/services/grpc/auth_interceptor.dart';
import 'package:mobile/src/services/grpc/grpc_flutter_client.dart';
import 'package:mobile/src/utils/retryutils.dart';
import 'package:riverpod/riverpod.dart';

class PetifiesService {
  AuthGatewayClient? _authClientInstance;
  Ref _ref;

  PetifiesService({
    required ref,
  }) : this._ref = ref;

  Future<AuthGatewayClient> get _authClient async {
    if (_authClientInstance == null) {
      _authClientInstance = AuthGatewayClient(
        await GrpcFlutterClient.getClient(),
        interceptors: [AuthInterceptor(ref: _ref)],
      );
    }
    return _authClientInstance!;
  }

  Future<PetifiesWithUserInfo> userCreatePetifies({
    required PetifiesType type,
    required String title,
    required String description,
    required String petName,
    required List<NetworkImageModel> images,
    required addressModel.Address address,
  }) async {
    UserCreatePetifiesRequest request = UserCreatePetifiesRequest(
      type: type,
      title: title,
      description: description,
      petName: petName,
      images: images.map((e) => Image(uri: e.uri, description: e.description)),
      address: Address(
        addressLineOne: address.addressLineOne,
        addressLineTwo: address.addressLineTwo,
        street: address.street,
        district: address.district,
        city: address.city,
        region: address.region,
        postalCode: address.postalCode,
        country: address.country,
        longitude: address.longitude,
        latitude: address.latitude,
      ),
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userCreatePetifies(request), _ref);
  }

  Future<UserPetifiesSession> userCreatePetifiesSession({
    required String petifiesId,
    required DateTime fromTime,
    required DateTime toTime,
  }) async {
    UserCreatePetifiesSessionRequest request = UserCreatePetifiesSessionRequest(
      petifiesId: petifiesId,
      fromTime: Timestamp.fromDateTime(fromTime),
      toTime: Timestamp.fromDateTime(toTime),
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userCreatePetifiesSession(request),
        _ref);
  }

  Future<PetifiesProposalWithUserInfo> userCreatePetifiesProposal({
    required String petifiesSessionId,
    required String proposal,
  }) async {
    UserCreatePetifiesProposalRequest request =
        UserCreatePetifiesProposalRequest(
            petifiesSessionId: petifiesSessionId, proposal: proposal);

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userCreatePetifiesProposal(request),
        _ref);
  }

  Future<ReviewWithUserInfo> userCreateReview({
    required String petifiesId,
    required String review,
    NetworkImageModel? image = null,
  }) async {
    UserCreateReviewRequest request = UserCreateReviewRequest(
      petifiesId: petifiesId,
      review: review,
      image: image != null
          ? Image(uri: image.uri, description: image.description)
          : null,
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userCreateReview(request), _ref);
  }

  Future<ListNearByPetifiesResponse> listNearByPetifies({
    required PetifiesType type,
    required double longitude,
    required double latitude,
    required double radius,
    int pageSize = 40,
    int offset = 0,
  }) async {
    ListNearByPetifiesRequest request = ListNearByPetifiesRequest(
      type: type,
      longitude: longitude,
      latitude: latitude,
      radius: radius,
      pageSize: pageSize,
      offset: offset,
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).listNearByPetifies(request), _ref);
  }

  Future<ListPetifiesByUserIdResponse> listPetifiesByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    ListPetifiesByUserIdRequest request = ListPetifiesByUserIdRequest(
        userId: userId, pageSize: pageSize, afterId: afterId);

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).listPetifiesByUserId(request), _ref);
  }

  Future<ListSessionsByPetifiesIdResponse> listSessionsByPetifiesId({
    required String petifiesId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    ListSessionsByPetifiesIdRequest request = ListSessionsByPetifiesIdRequest(
        petifiesId: petifiesId, pageSize: pageSize, afterId: afterId);

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).listSessionsByPetifiesId(request),
        _ref);
  }

  Future<ListProposalsBySessionIdResponse> listProposalsBySessionId({
    required String sessionId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    ListProposalsBySessionIdRequest request = ListProposalsBySessionIdRequest(
        sessionId: sessionId, pageSize: pageSize, afterId: afterId);

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).listProposalsBySessionId(request),
        _ref);
  }

  Future<ListProposalsByUserIdResponse> listProposalsByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    ListProposalsByUserIdRequest request = ListProposalsByUserIdRequest(
        userId: userId, pageSize: pageSize, afterId: afterId);

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).listProposalsByUserId(request), _ref);
  }

  Future<ListReviewsByPetifiesIdResponse> listReviewsByPetifiesId({
    required String petifiesId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    ListReviewsByPetifiesIdRequest request = ListReviewsByPetifiesIdRequest(
        petifiesId: petifiesId, pageSize: pageSize, afterId: afterId);

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).listReviewsByPetifiesId(request), _ref);
  }

  Future<ListReviewsByUserIdResponse> listReviewsByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    ListReviewsByUserIdRequest request = ListReviewsByUserIdRequest(
        userId: userId, pageSize: pageSize, afterId: afterId);

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).listReviewsByUserId(request), _ref);
  }
}
