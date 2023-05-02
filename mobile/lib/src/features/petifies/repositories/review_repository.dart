import 'dart:io';

import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/petifies_service.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/review.dart';

final reviewRepositoryProvider = Provider((ref) =>
    ReviewRepository(petifiesService: ref.read(petifiesServiceProvider)));

abstract class IReviewRepository {
  Future<ReviewModel> createReview({
    required String petifiesId,
    required String review,
    NetworkImageModel? image = null,
  });
  Future<List<ReviewModel>> listReviewsByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  });
  Future<List<ReviewModel>> listReviewsByPetifiesId({
    required String petifiesId,
    int pageSize = 40,
    String afterId = "",
  });
}

class ReviewRepository implements IReviewRepository {
  final PetifiesService _petifiesService;

  ReviewRepository({required PetifiesService petifiesService})
      : this._petifiesService = petifiesService;

  ReviewModel toReviewModel(ReviewWithUserInfo review) {
    return ReviewModel(
      id: review.id,
      author: BasicUserInfoModel(
        id: review.author.id,
        email: review.author.email,
        userAvatar: review.author.userAvatar,
        firstName: review.author.firstName,
        lastName: review.author.lastName,
      ),
      review: review.review,
      image: review.image.uri != ""
          ? NetworkImageModel(uri: review.image.uri)
          : null,
      petifiesId: review.petifiesId,
      createdAt: review.createdAt.toDateTime(),
    );
  }

  Future<ReviewModel> createReview({
    required String petifiesId,
    required String review,
    NetworkImageModel? image = null,
  }) async {
    final reviewResp = await _petifiesService.userCreateReview(
      petifiesId: petifiesId,
      review: review,
      image: image,
    );

    return toReviewModel(reviewResp);
  }

  Future<List<ReviewModel>> listReviewsByUserId({
    required String userId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    final reviews = await _petifiesService.listReviewsByUserId(
      userId: userId,
      pageSize: pageSize,
      afterId: afterId,
    );

    return reviews.reviews.map((e) => toReviewModel(e)).toList();
  }

  Future<List<ReviewModel>> listReviewsByPetifiesId({
    required String petifiesId,
    int pageSize = 40,
    String afterId = "",
  }) async {
    final reviews = await _petifiesService.listReviewsByPetifiesId(
      petifiesId: petifiesId,
      pageSize: pageSize,
      afterId: afterId,
    );

    return reviews.reviews.map((e) => toReviewModel(e)).toList();
  }
}
