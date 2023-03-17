import 'dart:async';
import 'dart:io';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/post/repository/file_repository.dart';
import 'package:mobile/src/models/video.dart';

final videoControllerProvider =
    AsyncNotifierProvider.autoDispose<VideoController, NetworkVideoModel?>(
        VideoController.new);

class VideoController extends AutoDisposeAsyncNotifier<NetworkVideoModel?> {
  @override
  FutureOr<NetworkVideoModel?> build() {
    return null;
  }

  Future<Either<Failure, NetworkVideoModel>> uploadVideo({
    required String uploaderID,
    required File video,
  }) async {
    final fileRepository = ref.read(fileRepositoryProvider);

    return fileRepository.uploadVideo(uploaderID: uploaderID, video: video);
  }

  Future<Either<Failure, void>> removeVideo({
    required String uri,
  }) async {
    final fileRepository = ref.read(fileRepositoryProvider);

    return fileRepository.removeFile(uri: uri);
  }
}
