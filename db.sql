USE [tsa]
GO

CREATE TABLE [dbo].[contact](
	[full_name] [varchar](512) NOT NULL,
	[email] [varchar](512) NOT NULL,
 CONSTRAINT [PK_contact] PRIMARY KEY CLUSTERED 
(
	[full_name] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO

CREATE TABLE [dbo].[contact_phone_number](
	[full_name] [varchar](512) NOT NULL,
	[number] [varchar](16) NOT NULL
) ON [PRIMARY]

GO

SET ANSI_PADDING OFF
GO

ALTER TABLE [dbo].[contact_phone_number]  WITH CHECK ADD  CONSTRAINT [FK_contact_phone_number_contact] FOREIGN KEY([full_name])
REFERENCES [dbo].[contact] ([full_name])
GO

ALTER TABLE [dbo].[contact_phone_number] CHECK CONSTRAINT [FK_contact_phone_number_contact]
GO

insert into tsa.dbo.contact (full_name, email) 
values 
	('Tsa Test', 'test@tsa.com.au'),
	('Chung Ho', 'chung@domain.com')

insert into tsa.dbo.contact_phone_number (full_name, number)
values
	('Tsa Test', '1800123456'),
	('Chung Ho', '0409123456'),
	('Chung Ho', '0860001234')
