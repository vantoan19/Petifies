import 'dart:io';

import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/post_service.dart';
import 'package:riverpod/riverpod.dart';

final postRepositoryProvider = Provider(
    (ref) => PostRepository(postService: ref.read(postServiceProvider)));

abstract class IPostRepository {
  Future<Either<Failure, PostModel>> createPost({
    required UserModel author,
    required String textContent,
    required List<NetworkImageModel> images,
    required List<NetworkVideoModel> videos,
  });
}

class PostRepository implements IPostRepository {
  final PostService _postService;

  PostRepository({required PostService postService})
      : this._postService = postService;

  Future<Either<Failure, PostModel>> createPost({
    required UserModel author,
    required String textContent,
    required List<NetworkImageModel> images,
    required List<NetworkVideoModel> videos,
  }) async {
    try {
      Post resp = await _postService.userCreatePost(
        textContent: textContent,
        images: images,
        videos: videos,
      );

      PostModel post = PostModel(
        owner: author,
        postActivity: "post",
        postTime: resp.createdAt,
        loveCount: 0,
        commentCount: 0,
      );
      return right(post);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }
}
