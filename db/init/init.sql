SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
CREATE SCHEMA IF NOT EXISTS `game_gacha` DEFAULT CHARACTER SET utf8mb4 ;
USE `game_gacha` ;

SET CHARSET utf8mb4;

-- -----------------------------------------------------
-- Table `game_gacha`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_gacha`.`users` (
    `id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
    `auth_token` VARCHAR(128) NOT NULL COMMENT '認証トークン',
    `name` VARCHAR(64) NOT NULL COMMENT 'ユーザ名',
    `high_score` INT UNSIGNED NOT NULL COMMENT 'ハイスコア',
    `coin` INT UNSIGNED NOT NULL COMMENT '所持コイン',
    `created_at` DATETIME NULL COMMENT '登録日時',
    `updated_at` DATETIME NULL COMMENT '最終プレイ日時',
    PRIMARY KEY (`id`),
    INDEX `idx_auth_token` (`auth_token` ASC)
)
ENGINE = InnoDB
COMMENT = 'ユーザ';


-- -----------------------------------------------------
-- Table `game_gacha`.`collection_items`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_gacha`.`collection_items` (
    `id` VARCHAR(128) NOT NULL COMMENT 'コレクションアイテムID',
    `name` VARCHAR(64) NOT NULL COMMENT 'コレクションアイテム名',
    `rarity` INT NOT NULL COMMENT 'レアリティ',
    `created_at` DATETIME NULL COMMENT '登録日時',
    PRIMARY KEY (`id`)
)
ENGINE = InnoDB
COMMENT = 'コレクションアイテム';


-- -----------------------------------------------------
-- Table `game_gacha`.`user_collection_item`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_gacha`.`user_collection_items` (
    `user_id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
    `collection_item_id` VARCHAR(128) NOT NULL COMMENT 'コレクションアイテムID',
    `created_at` DATETIME NULL COMMENT '登録日時',
    PRIMARY KEY (`user_id`, `collection_item_id`),
    INDEX `fk_user_collection_items_users_idx` (`user_id` ASC),
    INDEX `fk_user_collection_items_collection_items_idx` (`collection_item_id` ASC),
    CONSTRAINT `fk_user_collection_items_users`
    FOREIGN KEY (`user_id`)
    REFERENCES `game_gacha`.`users` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
    CONSTRAINT `fk_user_collection_items_collection_items`
    FOREIGN KEY (`collection_item_id`)
    REFERENCES `game_gacha`.`collection_items` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB
COMMENT = 'ユーザ所持コレクションアイテム';


-- -----------------------------------------------------
-- Table `game_gacha`.`gacha_probability`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_gacha`.`gacha_probabilities` (
    `collection_item_id` VARCHAR(128) NOT NULL COMMENT 'コレクションアイテムID',
    `ratio` INT UNSIGNED NOT NULL COMMENT '排出重み',
    INDEX `fk_gacha_probabilities_collection_items_idx` (`collection_item_id` ASC),
    PRIMARY KEY (`collection_item_id`),
    CONSTRAINT `fk_gacha_probabilities_collection_items_id`
    FOREIGN KEY (`collection_item_id`)
    REFERENCES `game_gacha`.`collection_items` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB
COMMENT = 'ガチャ排出情報';


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;