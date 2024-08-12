# Jarbis

Jarbis is an AI assistant inspired by J.A.R.V.I.S. from Iron Man, built using Go and leveraging various technologies such as Google Cloud Platform, FFmpeg, and speech recognition and synthesis APIs.

## Goal

The goal of Jarbis is to create a fully functional AI assistant that can understand and execute voice commands, engage in conversations, and assist with programming tasks. Jarbis aims to be a powerful tool for developers, allowing them to interact with their code repositories, run programs, and manipulate files through natural language commands.

In the future, Jarbis is envisioned to have the capability to control the entire operating system, providing a seamless and intuitive way to interact with computers using voice.

## Features

- Voice-based interaction: Jarbis uses speech recognition to understand voice commands and respond accordingly.
- Text-to-speech: Jarbis can convert its responses into natural-sounding speech using text-to-speech synthesis.
- Code repository manipulation: Jarbis can understand and execute commands related to code repositories, such as creating files, modifying code, and running programs.
- Intelligent assistance: Jarbis leverages AI technologies to provide intelligent responses and assist with programming tasks.

## Technologies Used

- Go programming language
- Google Cloud Platform
  - Google Cloud Speech-to-Text API for speech recognition
  - Google Cloud Text-to-Speech API for speech synthesis
- FFmpeg for audio recording and playback
- RobotGo for controlling mouse movements and keyboard inputs. 
- Git for version control and repository management

## Setup and Usage

1. Clone the repository:
   ```
   git clone https://github.com/EthanQuarry/jarbis.git
   ```

2. Set up the necessary dependencies and APIs:
   - Install Go and configure your Go environment
   - Set up a Google Cloud Platform project and enable the Speech-to-Text and Text-to-Speech APIs
   - Install FFmpeg

3. Configure the required environment variables:
   - `GOOGLE_APPLICATION_CREDENTIALS`: Path to your Google Cloud Platform credentials file
   - `GROQ_API_KEY`: API key for the Groq API (if applicable)
   - `INPUT_DEVICE_NAME`: Variable for your input device name

4. Build and run the Jarbis application:
   ```
   cd jarbis
   go build ./cmd/jarbis
   ./jarbis
   ```

5. Interact with Jarbis using voice commands or by typing in the console.

## Future Enhancements

- Integration with popular code editors and IDEs
- Support for multiple programming languages and frameworks
- Enhanced natural language understanding and context awareness
- Ability to control the operating system and perform system-level tasks
- Improved error handling and recovery mechanisms

## Contributing

Contributions to Jarbis are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

Jarbis is open-source software licensed under the [MIT License](LICENSE).
