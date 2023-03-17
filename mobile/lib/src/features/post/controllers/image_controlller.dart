import 'dart:async';
import 'dart:ffi';
import 'dart:io';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/post/repository/file_repository.dart';
import 'package:mobile/src/models/image.dart';

final imageControllerProvider =
    AsyncNotifierProvider.autoDispose<ImageController, NetworkImageModel?>(
        ImageController.new);

class ImageController extends AutoDisposeAsyncNotifier<NetworkImageModel?> {
  @override
  FutureOr<NetworkImageModel?> build() {
    return null;
  }

  Future<Either<Failure, NetworkImageModel>> uploadImage({
    required String uploaderID,
    required File image,
  }) async {
    final fileRepository = ref.read(fileRepositoryProvider);

    return fileRepository.uploadImage(uploaderID: uploaderID, image: image);
  }

  Future<Either<Failure, void>> removeImage({
    required String uri,
  }) async {
    final fileRepository = ref.read(fileRepositoryProvider);

    return fileRepository.removeFile(uri: uri);
  }
}
