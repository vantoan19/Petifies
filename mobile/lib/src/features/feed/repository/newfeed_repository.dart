import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:grpc/service_api.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/post_service.dart';

final newfeedRepositoryProvider = Provider<INewFeedRepository>(
    (ref) => NewFeedRepository(postService: ref.read(postServiceProvider)));

abstract class INewFeedRepository {
  Future<
      Tuple2<StreamController<ListNewFeedsRequest>,
          ResponseStream<ListNewFeedsResponse>>> getListNewfeedStream();
}

class NewFeedRepository implements INewFeedRepository {
  final PostService _postService;

  NewFeedRepository({
    required PostService postService,
  }) : _postService = postService;

  Future<
      Tuple2<StreamController<ListNewFeedsRequest>,
          ResponseStream<ListNewFeedsResponse>>> getListNewfeedStream() {
    return _postService.listNewFeeds();
  }
}
