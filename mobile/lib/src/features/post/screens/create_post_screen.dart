import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/post/controllers/image_controlller.dart';
import 'package:mobile/src/features/post/controllers/video_controller.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/providers/model_providers.dart';
import 'package:mobile/src/widgets/appbars/create_post_appbar.dart';
import 'package:image_picker/image_picker.dart';
import 'package:mobile/src/widgets/media_viewer/image_viewer.dart';
import 'package:mobile/src/widgets/pickers/visibility_picker.dart';
import 'package:video_player/video_player.dart';

class CreatePostScreen extends StatelessWidget {
  const CreatePostScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: CreatePostAppBar(),
      body: NewPostForm(),
      bottomNavigationBar: SizedBox(height: 45),
    );
  }
}

class NewPostForm extends ConsumerStatefulWidget {
  @override
  ConsumerState<ConsumerStatefulWidget> createState() => _NewPostFormState();
}

class _NewPostFormState extends ConsumerState<NewPostForm> {
  TextEditingController _textController = TextEditingController();
  List<File> _images = [];
  List<VideoPlayerController> _videoControllers = [];
  String? _visibility = "public";
  List<Future<Either<Failure, NetworkImageModel>>> _imagesFutures = [];
  List<Future<Either<Failure, NetworkVideoModel>>> _videosFutures = [];

  final ImagePicker _picker = ImagePicker();
  late FocusNode _textContentFocusNode;

  Future<void> _pickImages(String userID) async {
    final List<XFile>? images = await _picker.pickMultiImage();
    if (images != null) {
      setState(() {
        _images.addAll(images.map((img) => File(img.path)));
        _imagesFutures.addAll(
          images.map(
            (img) {
              return ref
                  .read(imageControllerProvider.notifier)
                  .uploadImage(uploaderID: userID, image: File(img.path));
            },
          ),
        );
      });
    }
    _textContentFocusNode.requestFocus();
  }

  Future<void> _pickVideo(String userID) async {
    final XFile? video = await _picker.pickVideo(source: ImageSource.gallery);
    if (video != null) {
      setState(() {
        _videoControllers.add(
          VideoPlayerController.file(File(video.path))
            ..initialize().then((_) {
              setState(() {});
            }),
        );

        _videosFutures.add(ref
            .read(videoControllerProvider.notifier)
            .uploadVideo(uploaderID: userID, video: File(video.path)));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  Future<void> _takePicture(String userID) async {
    final XFile? image = await _picker.pickImage(source: ImageSource.camera);
    if (image != null) {
      setState(() {
        _images.add(File(image.path));

        _imagesFutures.add(ref
            .read(imageControllerProvider.notifier)
            .uploadImage(uploaderID: userID, image: File(image.path)));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  Future<void> _takeVideo(String userID) async {
    final XFile? video = await _picker.pickVideo(source: ImageSource.camera);
    if (video != null) {
      setState(() {
        _videoControllers.add(
          VideoPlayerController.file(File(video.path))
            ..initialize().then((_) {
              setState(() {});
            }),
        );

        _videosFutures.add(ref
            .read(videoControllerProvider.notifier)
            .uploadVideo(uploaderID: userID, video: File(video.path)));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  void _removeImage(int idx) async {
    Future<Either<Failure, NetworkImageModel>> imgFuture = _imagesFutures[idx];
    setState(() {
      _images.removeAt(idx);
      _imagesFutures.removeAt(idx);
    });
    (await imgFuture).fold((l) => null, (r) {
      ref.read(imageControllerProvider.notifier).removeImage(uri: r.uri);
    });
  }

  @override
  void initState() {
    super.initState();
    _textContentFocusNode = FocusNode();
  }

  @override
  void dispose() {
    // Clean up the focus node when the Form is disposed.
    _textContentFocusNode.dispose();

    super.dispose();
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

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.max,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 0, 16),
            child: Row(
              children: [
                // User Avatar
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 0, 12, 0),
                  child: (userInfo.userAvatar != null)
                      ? CircleAvatar(
                          backgroundImage: NetworkImage(
                            userInfo.userAvatar!,
                          ),
                          backgroundColor: Colors.transparent,
                          radius: 24,
                        )
                      : CircleAvatar(
                          backgroundImage: AssetImage(
                            Constants.defaultAvatarPng,
                          ),
                          backgroundColor: Colors.transparent,
                          radius: 24,
                        ),
                ),
                // Visibility selection
                VisibilityPicker(
                  initialValue: _visibility,
                  onVisibilityChangedFunc: (value) => _visibility = value,
                ),
              ],
            ),
          ),
          // Post Content
          Expanded(
            child: SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  TextFormField(
                    autofocus: true,
                    controller: _textController,
                    decoration: InputDecoration(
                        hintText: 'What\'s happening?',
                        border: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        contentPadding: EdgeInsets.symmetric(vertical: 10)),
                    style: TextStyle(
                      fontSize: 20,
                    ),
                    maxLines: null,
                    focusNode: _textContentFocusNode,
                    keyboardType: TextInputType.multiline,
                  ),
                  SizedBox(
                    height: MediaQuery.of(context).size.width * 0.5,
                    child: ListView(
                      scrollDirection: Axis.horizontal,
                      children: [
                        ..._images.asMap().entries.map((entry) {
                          return ImageViewer(
                            image: entry.value,
                            uploaderID: userInfo.id,
                            removeAction: () => _removeImage(entry.key),
                          );
                        }).toList(),
                        ..._videoControllers.map(
                          (video) => AspectRatio(
                            aspectRatio: video.value.aspectRatio,
                            child: VideoPlayer(video),
                          ),
                        )
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  IconButton(
                    icon: Icon(Icons.camera_alt_outlined),
                    onPressed: () => _takePicture(userInfo.id),
                  ),
                  IconButton(
                    icon: Icon(Icons.videocam_outlined),
                    onPressed: () => _takeVideo(userInfo.id),
                  ),
                  IconButton(
                    icon: Icon(Icons.photo_library_outlined),
                    onPressed: () => _pickImages(userInfo.id),
                  ),
                  IconButton(
                    icon: Icon(Icons.video_library_outlined),
                    onPressed: () => _pickVideo(userInfo.id),
                  ),
                ],
              ),
              IconButton(
                icon: Icon(
                  Icons.keyboard_outlined,
                ),
                onPressed: () {
                  if (_textContentFocusNode.hasFocus) {
                    _textContentFocusNode.unfocus();
                  } else {
                    _textContentFocusNode.requestFocus();
                  }
                },
              ),
            ],
          ),
        ],
      ),
    );
  }
}
