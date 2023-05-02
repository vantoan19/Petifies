import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/petifies/repositories/review_repository.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/review.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'review_controllers.g.dart';

@riverpod
class CreateReviewController extends _$CreateReviewController {
  @override
  void build() {}

  Future<ReviewModel> createReview({
    required String petifiesId,
    required String review,
    Future<Either<Failure, NetworkImageModel>>? imageFuture = null,
  }) async {
    final reviewRepository = ref.read(reviewRepositoryProvider);

    NetworkImageModel? image = null;

    if (imageFuture != null) {
      final imageEither = await imageFuture;
      imageEither.fold((l) => null, (r) => image = r);
    }

    return reviewRepository.createReview(
        petifiesId: petifiesId, review: review, image: image);
  }
}

@riverpod
class ListReviewsByUserIdController extends _$ListReviewsByUserIdController {
  @override
  Future<List<ReviewModel>> build({required String userId}) async {
    final reviewRepository = ref.read(reviewRepositoryProvider);

    return reviewRepository.listReviewsByUserId(userId: userId);
  }

  void fetchMoreReviews() {
    state.whenData((value) async {
      state = AsyncLoading();
      state = await AsyncValue.guard(() async {
        final reviewRepository = ref.read(reviewRepositoryProvider);
        final reviews = await reviewRepository.listReviewsByUserId(
            userId: userId, afterId: value.last.id);
        return [...value, ...reviews];
      });
    });
  }
}

@riverpod
class ListReviewsByPetifidesIdController
    extends _$ListReviewsByPetifidesIdController {
  @override
  Future<List<ReviewModel>> build({required String petifiesId}) async {
    final reviewRepository = ref.read(reviewRepositoryProvider);

    return reviewRepository.listReviewsByPetifiesId(petifiesId: petifiesId);
  }

  void fetchMoreReviews() {
    state.whenData((value) async {
      state = AsyncLoading();
      state = await AsyncValue.guard(() async {
        final reviewRepository = ref.read(reviewRepositoryProvider);
        final reviews = await reviewRepository.listReviewsByPetifiesId(
            petifiesId: petifiesId, afterId: value.last.id);
        return [...value, ...reviews];
      });
    });
  }

  void addReview(ReviewModel review) {
    state.whenData((value) {
      state = AsyncData([...value, review]);
    });
  }
}
