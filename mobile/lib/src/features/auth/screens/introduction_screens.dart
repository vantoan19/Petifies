import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:introduction_screen/introduction_screen.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/custom_widgets/auth_appbar/auth_appbar.dart';

class IntroductionScreens extends StatelessWidget {
  const IntroductionScreens({super.key});

  List<PageViewModel> _getPages(width) {
    return [
      _IntroductionPage(
          title: "Join The Best Pet Lovers Community, Right Now.",
          bodyText: "Petfies is the best community for Pet Lovers where people"
              "can share momments, make friends, and "
              "especially find pet companions.",
          assetPath: Constants.firstIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: "Fall In With Pet Lovers Around The World.",
          bodyText: "With Petifies, you can easily find Pet Lovers"
              "around the world because Petifies is a huge network"
              "where you can find and add friends or join in groups.",
          assetPath: Constants.secondIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: "Let's Share Your Memories.",
          bodyText: "With Petifies, sharing your memories with your pets"
              "has never been as easy as now. Petfies allows you to "
              "post with videos and images, share stories within your network.",
          assetPath: Constants.thirdIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: "Keep In Touch With Your Friends.",
          bodyText: "In Petifies, Pet Lovers can send messages, comments on"
              "others' posts. Such a nice way to keep in touch with your friends.",
          assetPath: Constants.fourthIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: "Petifies, Fall In Love With Pets And Enjoy Your Happy Life"
              "With Animal Friends.",
          bodyText: "The most special thing that makes Petifies awesome is"
              "the system of Petifies. Petifies is a beautiful terms of"
              "cat/dog walking sessions, short-term/long-term cat/dog sittings,"
              "cat/dog adoptions.",
          assetPath: Constants.fifthIntroImagePath,
          deviceWidth: width),
    ];
  }

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.of(context).size.width;
    final GlobalKey<IntroductionScreenState> myKey = GlobalKey();

    final IntroductionScreen screen = IntroductionScreen(
      key: myKey,
      pages: _getPages(width),
      showSkipButton: true,
      showNextButton: true,
      done: const Text("Done"),
      skip: const Text("Skip"),
      next: const Text("Next"),
      onDone: () {
        Navigator.of(context).pushNamed("/signin");
      },
      dotsDecorator: DotsDecorator(
        size: const Size.square(10.0),
        activeSize: const Size(20.0, 10.0),
        activeColor: Theme.of(context).colorScheme.primary,
        color: Colors.black26,
        spacing: const EdgeInsets.symmetric(horizontal: 3.0),
        activeShape:
            RoundedRectangleBorder(borderRadius: BorderRadius.circular(25.0)),
      ),
    );

    return Scaffold(
      appBar: AuthAppBar(
        introScreenKey: myKey,
      ),
      body: screen,
    );
  }
}

class _IntroductionPage extends PageViewModel {
  static const headerStyle =
      TextStyle(fontSize: 34, fontWeight: FontWeight.w900);
  static const textStyle = TextStyle(fontSize: 16, fontWeight: FontWeight.w200);

  _IntroductionPage(
      {required title,
      required bodyText,
      required assetPath,
      required deviceWidth})
      : super(
            titleWidget: Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: Constants.horizontalScreenPadding),
              child: Text(
                title,
                style: headerStyle,
              ),
            ),
            bodyWidget: Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: Constants.horizontalScreenPadding),
              child: Text(
                bodyText,
                style: textStyle,
              ),
            ),
            image: SvgPicture.asset(
              assetPath,
              width: deviceWidth,
            ));
}
