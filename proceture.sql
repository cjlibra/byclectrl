DELIMITER $$
/*
select DISTINCT area_tb.area_name , moped_tb.moped_hphm ,type1_tb.dicword_wordname as typetype ,  color1_tb.dicword_wordname , 
	  owner_tb.owner_name , moped_tb.moped_id , tag_tb.tag_tagid ,moped_tb.moped_state , type2_tb.dicword_wordname as statestr
	  FROM owner_tb  JOIN moped_tb JOIN tag_tb   JOIN mopedowner_tb  
			ON moped_tb.moped_id = mopedowner_tb.moped_id AND  mopedowner_tb.owner_id = owner_tb.owner_id  
			JOIN mopedtag_tb ON mopedtag_tb.moped_id = moped_tb.moped_id AND mopedtag_tb.tag_id = tag_tb.tag_id  
			JOIN area_tb ON area_tb.area_id = moped_tb.area_id  
			JOIN  dicword_tb  AS type2_tb  ON  type2_tb.dicword_dictypeid = 8 AND tag_tb.tag_state = type2_tb.dicword_wordid  
			JOIN  dicword_tb  AS type1_tb  ON  type1_tb.dicword_dictypeid = 6 AND moped_tb.moped_type = type1_tb.dicword_wordid 
			JOIN   dicword_tb  AS color1_tb  ON   color1_tb.dicword_dictypeid = 7
			 AND moped_tb.moped_colorid = color1_tb.dicword_wordid  
			WHERE mopedtag_tb.mopedtag_state = 1 and moped_tb.moped_hphm like "%%%s%%" and owner_tb.owner_name like "%%%s%%" and (owner_tb.owner_state = 1) order by owner_tb.owner_id 
			
			*/
	
CREATE
    /*[DEFINER = { user | CURRENT_USER }]*/
    PROCEDURE `mopedmanage`.`getMopedBynameOrHphm`(IN v_hphm VARCHAR(45) ,  IN v_ownername VARCHAR(45))
    /*LANGUAGE SQL
    | [NOT] DETERMINISTIC
    | { CONTAINS SQL | NO SQL | READS SQL DATA | MODIFIES SQL DATA }
    | SQL SECURITY { DEFINER | INVOKER }
    | COMMENT 'string'*/
    BEGIN
      
        DROP TABLE IF EXISTS tmp_out_tb;
        CREATE	TEMPORARY TABLE tmp_out_tb (  \
			area_name  VARCHAR(45),\
			moped_hphm VARCHAR(45),\
			typetype VARCHAR(45),
			dicword_wordname VARCHAR(45),
			owner_name VARCHAR(45),
			moped_id  INTEGER ,
			tag_tagid VARCHAR(45),
			moped_state  INTEGER ,
			statestr VARCHAR(45)
			)  	;
		BEGIN 
			DECLARE s_Areaname VARCHAR(45);
			DECLARE s_Hphm VARCHAR(45);
			DECLARE s_Typetype VARCHAR(45);
			DECLARE s_Color VARCHAR(45);
			DECLARE s_Ownername VARCHAR(45);
			DECLARE i_Moped_id INT(10) ;
			DECLARE i_Moped_state INT(10);
			DECLARE s_Tag_tagid VARCHAR(45);
			DECLARE s_tag_state VARCHAR(45);
			
			DECLARE i_moped_type INT(10);
			DECLARE i_areaid INT(10);
			DECLARE i_colorid INT(10);
			DECLARE i_ownerid INT(10);
			DECLARE i_Tag_tagid INT(10);
			DECLARE i_tag_state INT(10);
		   
			DECLARE done INT;
			
		   
		   
			DECLARE p_moped_tb CURSOR FOR SELECT moped_id , moped_hphm , moped_type , area_id ,moped_state , moped_colorid FROM moped_tb WHERE moped_hphm LIKE    CONCAT('%',v_hphm,'%');
			
			 DECLARE CONTINUE HANDLER FOR NOT FOUND SET done=1;
			
			OPEN p_moped_tb  ;
			cursor_loop_moped : LOOP	
						
				FETCH p_moped_tb INTO 	i_moped_id , s_Hphm , i_moped_type , i_areaid , i_Moped_state , i_colorid ;
				BEGIN 
					DECLARE p_mopedowner_tb CURSOR FOR SELECT owner_id FROM mopedowner_tb WHERE moped_id = i_moped_id  ;
					OPEN p_mopedowner_tb ;
					cursor_loop_mopedowner : LOOP
						FETCH p_mopedowner_tb INTO   i_ownerid ;
						BEGIN
							DECLARE p_owner_tb CURSOR FOR SELECT owner_name FROM owner_tb WHERE ownerid =  i_ownerid  AND owner_name LIKE CONCAT('%', v_ownername,'%') AND ( owner_state = 1) ORDER BY  owner_id;
							OPEN p_owner_tb ;
							cursor_loop_owner : LOOP
								FETCH p_owner_tb INTO s_Ownername ;
								BEGIN 
									DECLARE p_mopedtag_tb CURSOR FOR SELECT tag_id FROM mopedtag_tb WHERE moped_id = i_moped_id ;
									OPEN p_mopedtag_tb;
									cursor_loop_mopedtag : LOOP
									
										FETCH p_mopedtag_tb INTO i_Tag_tagid ;
										
										SELECT tag_tagid,tag_state INTO s_Tag_tagid ,i_tag_state FROM tag_tb WHERE tag_id = i_Tag_tagid ;
										SELECT area_name INTO s_Areaname FROM area_tb WHERE area_id = i_areaid ;
										SELECT dicword_wordname INTO s_Color FROM dicword_tb WHERE  dicword_dictypeid = 7 AND dicword_wordid =  i_colorid;
										SELECT dicword_wordname INTO s_Typetype FROM dicword_tb WHERE  dicword_dictypeid = 6 AND dicword_wordid =  i_moped_type;
										SELECT dicword_wordname INTO s_tag_state FROM dicword_tb WHERE  dicword_dictypeid = 8 AND dicword_wordid =  i_tag_state;
										INSERT INTO tmp_out_tb VALUES(s_Areaname ,s_Hphm ,s_Typetype , s_Color ,s_Ownername ,i_moped_id ,s_Tag_tagid ,i_Moped_state , s_tag_state ) ;
										
										
										IF done=1 THEN
											SET done = 0;
											LEAVE cursor_loop_mopedtag;
										END IF;
									
									END LOOP cursor_loop_mopedtag;
									CLOSE p_mopedtag_tb ;
								
									IF done=1 THEN
										SET done = 0;
										LEAVE cursor_loop_owner;
									END IF;
								END;
							END LOOP cursor_loop_owner ;
							CLOSE p_owner_tb;
							
						
							IF done=1 THEN
								SET done = 0;
								LEAVE cursor_loop_mopedowner;
							END IF;
					    END;
					END LOOP cursor_loop_mopedowner ;
					CLOSE p_mopedowner_tb ;
					
			   
					IF done=1 THEN
						SET done = 0;
						LEAVE cursor_loop_moped;
					END IF;
				END;	
			END LOOP cursor_loop_moped ;	
			CLOSE p_moped_tb ;
		 END;
	
	SELECT DISTINCT * FROM tmp_out_tb ;
	 
    END$$

DELIMITER ;