import 'dart:async';
import 'dart:io';

import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/media/repositores/file_repository.dart';
import 'package:mobile/src/models/image.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'image_controller.g.dart';

@Riverpod(keepAlive: false)
class ImageController extends _$ImageController {
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
