import 'dart:io';

import 'package:fpdart/fpdart.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/post_service.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

final postRepositoryProvider = Provider(
    (ref) => PostRepository(postService: ref.read(postServiceProvider)));

abstract class IPostRepository {
  Future<Either<Failure, PostModel>> createPost({
    required BasicUserInfoModel author,
    String visibility = "public",
    required String activity,
    required String textContent,
    required List<NetworkImageModel> images,
    required List<NetworkVideoModel> videos,
  });
  Future<ResponseStream<StreamLoveCountResponse>> getPostLoveCountStream({
    required String postID,
  });
  Future<ResponseStream<StreamCommentCountResponse>> getPostCommentCountStream({
    required String postID,
  });
  Future<Either<Failure, bool>> userToggleLoveReactPost({
    required String postID,
  });
}

class PostRepository implements IPostRepository {
  final PostService _postService;

  PostRepository({required PostService postService})
      : this._postService = postService;

  Future<Either<Failure, PostModel>> createPost({
    required BasicUserInfoModel author,
    String visibility = "public",
    required String activity,
    required String textContent,
    required List<NetworkImageModel> images,
    required List<NetworkVideoModel> videos,
  }) async {
    try {
      PostWithUserInfo resp = await _postService.userCreatePost(
        visibility: visibility,
        activity: activity,
        textContent: textContent,
        images: images,
        videos: videos,
      );

      PostModel post = PostModel(
        id: resp.id,
        owner: author,
        postActivity: "post",
        createdAt: resp.createdAt.toDateTime(),
        textContent: resp.content == "" ? null : resp.content,
        images: resp.images
            .map((image) => NetworkImageModel(uri: image.uri))
            .toList(),
        videos: resp.videos
            .map((video) => NetworkVideoModel(uri: video.uri))
            .toList(),
        hasReacted: resp.hasReacted,
        loveCount: 0,
        commentCount: 0,
      );
      return right(post);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<ResponseStream<StreamLoveCountResponse>> getPostLoveCountStream({
    required String postID,
  }) {
    return _postService.getLoveCount(postID: postID);
  }

  Future<ResponseStream<StreamCommentCountResponse>> getPostCommentCountStream({
    required String postID,
  }) {
    return _postService.getCommentCount(postID: postID);
  }

  Future<Either<Failure, bool>> userToggleLoveReactPost({
    required String postID,
  }) async {
    try {
      UserToggleLoveResponse resp =
          await _postService.userToggleLoveReactPost(postID: postID);

      return right(resp.hasReacted.value);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }
}
