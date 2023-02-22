import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:introduction_screen/introduction_screen.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/constants/languages.dart';
import 'package:mobile/src/widgets/auth_appbar/auth_appbar.dart';

class IntroductionScreens extends StatelessWidget {
  const IntroductionScreens({super.key});

  List<PageViewModel> _getPages(context, width) {
    return [
      _IntroductionPage(
          title: Language.translate(context).firstIntroTitle,
          bodyText: Language.translate(context).firstIntroText,
          assetPath: Constants.firstIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: Language.translate(context).secondIntroTitle,
          bodyText: Language.translate(context).secondIntroText,
          assetPath: Constants.secondIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: Language.translate(context).thirdIntroTitle,
          bodyText: Language.translate(context).thirdIntroText,
          assetPath: Constants.thirdIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: Language.translate(context).fourthIntroTitle,
          bodyText: Language.translate(context).fourthIntroText,
          assetPath: Constants.fourthIntroImagePath,
          deviceWidth: width),
      _IntroductionPage(
          title: Language.translate(context).fifthIntroTitle,
          bodyText: Language.translate(context).fifthIntroText,
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
      pages: _getPages(context, width),
      showSkipButton: true,
      showNextButton: true,
      done: Text(Language.translate(context).doneLabel),
      skip: Text(Language.translate(context).skipLabel),
      next: Text(Language.translate(context).nextLabel),
      onDone: () {
        Navigator.of(context).pushNamed("/signin");
      },
      dotsDecorator: DotsDecorator(
        size: const Size.square(10.0),
        activeSize: const Size(20.0, 10.0),
        activeColor: Theme.of(context).colorScheme.primary,
        color: Theme.of(context).colorScheme.inversePrimary,
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
          ),
        );
}
