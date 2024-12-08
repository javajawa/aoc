knownAntennae: dict[bytes, set[complex]] = {}
knownAntiNodes: set[complex] = set()

with open("2024/input-8.txt", "rb") as input_stream:
    data = input_stream.readlines()

nonFrequencies = frozenset({46, 10})  # b'.' and b'\n'
size = len(data)

for row, line in enumerate(data):
    for col, frequency in enumerate(line):
        if frequency in nonFrequencies:
            continue

        knownAntennae.setdefault(frequency, set())

        newAntennaPoint = complex(col, row)

        for previousAntenna in knownAntennae[frequency]:
            # print("Calculating lines based from", previousAntenna, newAntennaPoint)

            vector = previousAntenna - newAntennaPoint
            # print(" ", "Vector is", vector)

            newAntiNode = previousAntenna
            while newAntiNode.real >= 0 and newAntiNode.imag >= 0 and newAntiNode.real < size and newAntiNode.imag < size:
                # print(" ", "Adding node", newAntiNode)
                knownAntiNodes.add(newAntiNode)
                newAntiNode += vector

            newAntiNode = newAntennaPoint
            while newAntiNode.real >= 0 and newAntiNode.imag >= 0 and newAntiNode.real < size and newAntiNode.imag < size:
                # print(" ", "Adding node", newAntiNode)
                knownAntiNodes.add(newAntiNode)
                newAntiNode -= vector

        knownAntennae[frequency].add(newAntennaPoint)

print(len(knownAntiNodes))
