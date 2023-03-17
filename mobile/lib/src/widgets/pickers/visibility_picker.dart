// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';

class VisibilityPicker extends StatefulWidget {
  final String? initialValue;
  final Function(String?) onVisibilityChangedFunc;

  const VisibilityPicker({
    Key? key,
    required this.initialValue,
    required this.onVisibilityChangedFunc,
  }) : super(key: key);

  @override
  State<VisibilityPicker> createState() => _VisibilityPickerState();
}

class _VisibilityPickerState extends State<VisibilityPicker> {
  String? _currentVisibility;

  @override
  void initState() {
    super.initState();
    _currentVisibility = widget.initialValue;
  }

  void _onVisibilityChange(String value, BuildContext context) {
    NavigatorUtil.goBack(context);
    setState(() {
      _currentVisibility = value;
      widget.onVisibilityChangedFunc(value);
    });
  }

  Icon _getVisibilityIcon() {
    if (_currentVisibility == null) {
      return Icon(
        size: 22,
        Icons.visibility,
        color: Themes.blueColor,
      );
    } else if (_currentVisibility == "public") {
      return Icon(
        size: 22,
        Icons.public,
        color: Themes.blueColor,
      );
    } else if (_currentVisibility == "friends") {
      return Icon(
        size: 22,
        Icons.group,
        color: Themes.blueColor,
      );
    } else {
      return Icon(
        size: 22,
        Icons.lock,
        color: Themes.blueColor,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      style: ElevatedButton.styleFrom(
        shape: RoundedRectangleBorder(
          side: BorderSide(color: Themes.blueColor),
          borderRadius: BorderRadius.circular(18.0),
        ),
        backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 4, 0),
            child: _getVisibilityIcon(),
          ),
          Text(
            (_currentVisibility != null) ? _currentVisibility! : "visibility",
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w900,
              color: Themes.blueColor,
            ),
          ),
          Icon(
            size: 22,
            Icons.expand_more,
            color: Themes.blueColor,
          ),
        ],
      ),
      onPressed: () {
        showModalBottomSheet<void>(
          context: context,
          enableDrag: false,
          shape: const RoundedRectangleBorder(
            borderRadius: BorderRadius.only(
              topLeft: Radius.circular(24),
              topRight: Radius.circular(24),
            ),
          ),
          barrierColor: Themes.lightGreyColorTransparent,
          backgroundColor: Theme.of(context).colorScheme.tertiary,
          builder: (BuildContext context) {
            return SizedBox(
              height: MediaQuery.of(context).size.height * 0.6,
              child: Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.start,
                  mainAxisSize: MainAxisSize.max,
                  children: <Widget>[
                    Container(
                      width: 60,
                      height: 5,
                      margin: const EdgeInsets.symmetric(vertical: 15),
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.all(Radius.circular(10)),
                        color: Theme.of(context).colorScheme.secondary,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 25,
                        vertical: 15,
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            "Who can see your post?",
                            style: TextStyle(
                              fontSize: 22,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.symmetric(vertical: 15),
                            child: Text(
                              "Your post will be visible on the feed, on your profile and in search results",
                              style: TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.w300,
                                color: Theme.of(context).colorScheme.secondary,
                              ),
                            ),
                          ),
                          _VisibilityOption(
                            optionLabel: "Anyone",
                            description: "Anyone on or off Petifies",
                            icon: Icon(Icons.public),
                            action: () =>
                                _onVisibilityChange("public", context),
                            isSelected: (_currentVisibility == "public"),
                          ),
                          _VisibilityOption(
                            optionLabel: "Friends",
                            description: "Friends on Petifies",
                            icon: Icon(Icons.group),
                            action: () =>
                                _onVisibilityChange("friends", context),
                            isSelected: (_currentVisibility == "friends"),
                          ),
                          _VisibilityOption(
                            optionLabel: "Only me",
                            description: "Only you",
                            icon: Icon(Icons.lock),
                            action: () =>
                                _onVisibilityChange("private", context),
                            isSelected: (_currentVisibility == "private"),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }
}

class _VisibilityOption extends StatelessWidget {
  final Widget icon;
  final String optionLabel;
  final String description;
  final VoidCallback action;
  final bool isSelected;

  const _VisibilityOption({
    Key? key,
    required this.icon,
    required this.optionLabel,
    required this.description,
    required this.action,
    required this.isSelected,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTap: action,
      child: Row(
        mainAxisSize: MainAxisSize.max,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 20, 15, 20),
            child: icon,
          ),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  optionLabel,
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Text(
                  description,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w200,
                  ),
                ),
              ],
            ),
          ),
          if (isSelected)
            Icon(
              Icons.check_circle,
              color: Themes.blueColor,
            ),
          if (!isSelected) Icon(Icons.radio_button_unchecked),
        ],
      ),
    );
  }
}
