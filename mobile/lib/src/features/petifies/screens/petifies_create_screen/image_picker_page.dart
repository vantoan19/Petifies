import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:image_picker/image_picker.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/media/controllers/image_controller.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_create_screen.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/widgets/images/image_upload_viewer.dart';

class ImagePickerPage extends ConsumerStatefulWidget {
  const ImagePickerPage({super.key});

  @override
  ConsumerState<ImagePickerPage> createState() => _ImagePickerPageState();
}

class _ImagePickerPageState extends ConsumerState<ImagePickerPage> {
  final ImagePicker _picker = ImagePicker();

  Future<void> _pickImages(String userID) async {
    final List<XFile>? images = await _picker.pickMultiImage();
    if (images != null) {
      setState(() {
        ref
            .read(imageFilesProvider.notifier)
            .addImages(images.map((img) => File(img.path)).toList());
        ref.read(imageFuturesProvider.notifier).addImages(
              images.map(
                (img) {
                  return ref
                      .read(imageControllerProvider.notifier)
                      .uploadImage(uploaderID: userID, image: File(img.path));
                },
              ).toList(),
            );
      });
    }
  }

  Future<void> _takePicture(String userID) async {
    final XFile? image = await _picker.pickImage(source: ImageSource.camera);
    if (image != null) {
      setState(() {
        ref.read(imageFilesProvider.notifier).addImages([File(image.path)]);
        ref.read(imageFuturesProvider.notifier).addImages([
          ref
              .read(imageControllerProvider.notifier)
              .uploadImage(uploaderID: userID, image: File(image.path))
        ]);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final user = ref.watch(myUserProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    if (userInfo == null) {
      return Placeholder();
    }

    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(
                Constants.petifiesExpoloreHorizontalPadding,
                28,
                Constants.petifiesExpoloreHorizontalPadding,
                20,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    "Add some photos of your pet",
                    style: TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(top: 8.0),
                    child: Text(
                      "You'll need to upload at least 3 photos to get started. You can add more or make changes later. First image will be choosed as cover photo.",
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w300,
                        color: Theme.of(context).colorScheme.secondary,
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(top: 24),
                    child: Row(
                      children: [
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: () {
                              _pickImages(userInfo.id);
                            },
                            icon: Icon(
                              Icons.add,
                              color:
                                  Theme.of(context).colorScheme.inversePrimary,
                            ),
                            label: Text(
                              "Add photos",
                              style: TextStyle(
                                color: Theme.of(context)
                                    .colorScheme
                                    .inversePrimary,
                              ),
                            ),
                            style: OutlinedButton.styleFrom(
                              alignment: Alignment.centerLeft,
                              padding: EdgeInsets.all(16),
                              shape: RoundedRectangleBorder(
                                  borderRadius:
                                      BorderRadius.all(Radius.circular(12))),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(top: 8.0),
                    child: Row(
                      children: [
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: () {
                              _takePicture(userInfo.id);
                            },
                            icon: Icon(
                              Icons.camera_alt_outlined,
                              color:
                                  Theme.of(context).colorScheme.inversePrimary,
                            ),
                            label: Text(
                              "Take new photos",
                              style: TextStyle(
                                color: Theme.of(context)
                                    .colorScheme
                                    .inversePrimary,
                              ),
                            ),
                            style: OutlinedButton.styleFrom(
                              alignment: Alignment.centerLeft,
                              padding: EdgeInsets.all(16),
                              shape: RoundedRectangleBorder(
                                  borderRadius:
                                      BorderRadius.all(Radius.circular(12))),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(top: 8.0),
                    child: const _ImagesPreviewer(),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ImagesPreviewer extends ConsumerWidget {
  const _ImagesPreviewer({super.key});

  void _removeImage(int idx, WidgetRef ref) {
    ref.read(imageFilesProvider.notifier).removeImage(idx);
    ref.read(imageFuturesProvider.notifier).removeImage(idx);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    if (userInfo == null) {
      return Placeholder();
    }

    final images = ref.watch(imageFilesProvider);

    return SizedBox(
      height: MediaQuery.of(context).size.width * 0.5,
      child: ListView(
        scrollDirection: Axis.horizontal,
        children: [
          ...images.asMap().entries.map((entry) {
            return ImageUploadViewer(
              image: entry.value,
              uploaderID: userInfo.id,
              removeAction: () => _removeImage(entry.key, ref),
            );
          }).toList(),
        ],
      ),
    );
  }
}
