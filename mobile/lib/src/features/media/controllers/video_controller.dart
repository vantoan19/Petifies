import 'dart:async';
import 'dart:io';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/media/repositores/file_repository.dart';
import 'package:mobile/src/models/video.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'video_controller.g.dart';

@Riverpod(keepAlive: false)
class VideoController extends _$VideoController {
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
