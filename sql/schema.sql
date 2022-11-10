/**
Cell Types
Type - "PGN"
**/
CREATE TABLE CellTypes (
    Type varchar(255) primary key NOT NULL
);

/**
Genes
Name - "Plp1"
Chr - "X"
Start - 1.37E8
Stop  - 1.37E8
**/
CREATE TABLE Genes (
    Name varchar(255) primary key,
    Chr varchar(3) NOT NULL,
    Start int NOT NULL,
    End int NOT NULL,
    Check (Start < End)
);

/**
Loci
ID: "chrX:305000000-3100000000"
Chr: "X"
Start: 305000000
End:   310000000
**/
CREATE TABLE Loci (
    ID varchar(255) primary key,
    Chr varchar(3) NOT NULL,
    Start INT NOT NULL,
    End INT NOT NULL,
    CHECK (Start < End)
);

/**
Motif Models
Name: "BHE40_MOUSE.H11MO.0A"
Length: 9
Quality: "A"
Uniprot ID: "BHE40_MOUSE"
Transcription Factor: "MOUSE:Bhlhe40"
TF Family: "Hair-related factors{1.2.4}"
Entrez Gene: 20893
**/
CREATE TABLE MotifModels (
    Name varchar(255) primary key,
    Length int NOT NULL,
    Quality char(1) NOT NULL,
    UniprotID varchar(255),
    TranscriptionFactor varchar(255) NOT NULL,
    TFFamily varchar(255) NOT NULL,
    EntrezGene int NOT NULL,
    CHECK (Length > 0)
);

/**
Motif Instance
Cell Type - "PGN"
Chr - "X"
Start - 250000
Forward - True
Threshold Score - 11.0796049119
Locus ID - "ChrX:3500000-35050000"
Model - "PBX1_MOUSE.H11MO.2.C"
**/
CREATE TABLE MotifInstances (
    CellType varchar(255) NOT NULL,
    Chr varchar(3) NOT NULL,
    Start INT NOT NULL,
    Forward BOOLEAN NOT NULL,
    ThresholdScore FLOAT NOT NULL,
    LocusID varchar(255) NOT NULL,
    Model varchar(255) NOT NULL,
    PRIMARY KEY (CellType, Chr, Start),
    FOREIGN KEY (CellType) REFERENCES CellTypes(Type),
    FOREIGN KEY (LocusID) REFERENCES Loci(ID),
    FOREIGN KEY (Model) REFERENCES MotifModels(Name)
);

/**
Interactions
Cell Type - "PGN"
ID - Unique ID (Sqllite recommends just leaving as int PRIMARY KEY)
It'll auto fill UID if not defined.
**/
CREATE TABLE Interactions (
    CellType varchar(255) NOT NULL,
    ID int PRIMARY KEY,
    FOREIGN KEY (CellType) REFERENCES CellTypes(Type)
);

/**
Gene expression
Cell Type - "PGN"
Gene - "Plp1"
Expression Level - 11.59034242
**/
CREATE TABLE GeneExpression (
    CellType varchar(255),
    Gene varchar(255),
    ExpressionLevel FLOAT NOT NULL,
    FOREIGN KEY (CellType) REFERENCES CellTypes(Type),
    FOREIGN KEY (Gene) REFERENCES Genes(Name),
    PRIMARY KEY (CellType, Gene)
);

/**
Interaction Participation
Locus - "chrX:350000-3550000"
Interaction ID - 1234
**/
CREATE TABLE InteractionParticipation (
    Locus varchar(255),
    Interaction varchar(255),
    FOREIGN KEY (Locus) REFERENCES Loci(ID),
    FOREIGN KEY (Interaction) REFERENCES Interactions(ID),
    PRIMARY KEY (Locus, Interaction)
);

/**
Gene Located in Locus
Locus - "chrX:350000-3550000"
Gene - "Plp1"
**/
CREATE TABLE GeneInLocus (
    Locus varchar(255),
    Gene varchar(255),
    FOREIGN KEY (Locus) REFERENCES Loci(ID),
    FOREIGN KEY (Gene) REFERENCES Genes(Name),
    PRIMARY KEY (Locus, Gene)
);