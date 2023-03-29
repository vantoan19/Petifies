// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/comment/screens/create_comment_screen.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/posts/post.dart';

class ReplyTextField extends ConsumerStatefulWidget {
  final bool autoFocus;
  const ReplyTextField({
    Key? key,
    required this.autoFocus,
  }) : super(key: key);

  @override
  ConsumerState<ReplyTextField> createState() => ReplyTextFieldState();
}

class ReplyTextFieldState extends ConsumerState<ReplyTextField> {
  late FocusNode _myFocusNode;

  @override
  void initState() {
    super.initState();

    _myFocusNode = FocusNode();
    _myFocusNode.addListener(() {
      if (_myFocusNode.hasFocus) {
        _myFocusNode.unfocus();

        showModalBottomSheet(
            context: context,
            isScrollControlled: true,
            useSafeArea: true,
            barrierColor: Theme.of(context).scaffoldBackgroundColor,
            builder: (context) {
              final isPostTarget = ref.read(isPostContextProvider);
              return ProviderScope(
                overrides: [
                  isPostContextProvider.overrideWithValue(isPostTarget),
                  postInfoProvider
                      .overrideWithValue(ref.read(postInfoProvider)),
                  if (!isPostTarget)
                    commentInfoProvider
                        .overrideWithValue(ref.read(commentInfoProvider)),
                ],
                child: CreateCommentScreen(),
              );
            });
      }
    });
  }

  @override
  void dispose() {
    _myFocusNode.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 20),
      child: TextField(
        autofocus: widget.autoFocus,
        focusNode: _myFocusNode,
        decoration: InputDecoration(
            filled: true,
            hintText: 'Write your reply',
            contentPadding:
                EdgeInsets.symmetric(vertical: 10.0, horizontal: 20.0),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(24.0),
              borderSide: BorderSide.none,
            ),
            constraints: BoxConstraints(maxHeight: 40)),
      ),
    );
    ;
  }
}
